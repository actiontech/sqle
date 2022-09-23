package cloudbeaver_wrapper

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/controller"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/resolver"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/log"
	sqleModel "github.com/actiontech/sqle/sqle/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type gqlBehavior struct {
	useLocalHandler     bool
	needModifyRemoteRes bool
	disable             bool
	// 预处理主要用于在真正使用前处理前端传递的参数, 比如需要接收int, 但收到float, 则可以在此处调整参数类型
	preprocessing func(ctx echo.Context, params *graphql.RawParams) error
}

var gqlHandlerRouters = map[string] /* gql operation name */ gqlBehavior{
	"authLogin": {
		useLocalHandler:     true,
		needModifyRemoteRes: true,
		preprocessing: func(ctx echo.Context, params *graphql.RawParams) error {
			// 还原参数中的用户名
			if credentials, ok := params.Variables["credentials"].(map[string]interface{}); ok {
				if credentials["user"] != nil {
					params.Variables["credentials"].(map[string]interface{})["user"] = service.GenerateCloudBeaverUserName(fmt.Sprintf("%v", params.Variables["credentials"].(map[string]interface{})["user"]))
				}
			}

			// 更新context
			body, err := json.Marshal(params)
			if err != nil {
				return err
			}
			ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
			ctx.Request().ContentLength = int64(len(body))
			return nil
		},
	},
	"asyncSqlExecuteQuery": {
		useLocalHandler:     true,
		needModifyRemoteRes: false,
		preprocessing: func(ctx echo.Context, params *graphql.RawParams) (err error) {
			// json中没有int类型, 这将导致执行json.Unmarshal()时int会被当作float64, 从而导致后面出现类型错误的异常
			if filter, ok := params.Variables["filter"].(map[string]interface{}); ok {
				if filter["limit"] != nil {
					params.Variables["filter"].(map[string]interface{})["limit"], err = strconv.Atoi(fmt.Sprintf("%v", params.Variables["filter"].(map[string]interface{})["limit"]))
				}
			}
			return err
		},
	},
	"getActiveUser": {
		useLocalHandler:     true,
		needModifyRemoteRes: true,
	},
	"configureServer": {
		disable: true,
	}, "createUser": {
		disable: true,
	}, "setUserCredentials": {
		disable: true,
	}, "enableUser": {
		disable: true,
	}, "grantUserRole": {
		disable: true,
	}, "setConnections": {
		disable: true,
	}, "saveUserMetaParameters": {
		disable: true,
	}, "deleteUser": {
		disable: true,
	}, "createRole": {
		disable: true,
	}, "updateRole": {
		disable: true,
	}, "deleteRole": {
		disable: true,
	}, "authChangeLocalPassword": {
		disable: true,
	},
}

func StartApp(e *echo.Echo) {
	if !service.IsCloudBeaverConfigured() {
		return
	}
	fmt.Println("cloudbeaver wrapper is configured")

	cfg := service.GetSQLQueryConfig()
	protocol := "http"
	if cfg.EnableHttps {
		protocol = "https"
	}
	url2, err := url.Parse(fmt.Sprintf("%v://%v:%v", protocol, cfg.CloudBeaverHost, cfg.CloudBeaverPort))
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url2,
		},
	}
	q := e.Group(service.CbRootUri)

	q.Use(RedirectCookie())
	q.Use(GraphqlDistributor())
	q.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Skipper:  middleware.DefaultSkipper,
		Balancer: middleware.NewRandomBalancer(targets),
	}))
}

// login页面无法访问sql_query页面的cookie, 这将导致登录SQLE时无法判断CloudBeaver当前登陆状态, 所以需要将cookie放到根目录下, 使用时在还原回原来的位置
func RedirectCookie() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cookie string
			// 如果用户传了cookie, 则还原cookie的路径
			for _, c := range c.Cookies() {
				if c.Name == "cb-session-id-sqle" {
					cookie = c.Value
				}
			}
			c.Request().Header.Del("Cookie")
			c.Request().Header.Set("Cookie", "cb-session-id="+cookie)

			return next(c)
		}
	}
}

// 登录CloudBeaver不应该影响SQLE登录
func TriggerLogin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resWrite := &responseProcessWriter{tmp: &bytes.Buffer{}, ResponseWriter: c.Response().Writer}
			c.Response().Writer = resWrite
			respFunc := func() error {
				_, err := resWrite.ResponseWriter.Write(resWrite.tmp.Bytes())
				return err
			}

			err := next(c)
			if err != nil {
				return err
			}

			// 如果登陆失败, userName应该取不出来
			userName, ok := c.Get(config.LoginUserNameKey).(string)
			if !ok || !service.IsCloudBeaverConfigured() {
				return respFunc()
			}

			l := log.NewEntry()
			s := sqleModel.GetStorage()
			user, _, err := s.GetUserByName(userName)
			if err != nil {
				l.Errorf("get user info err: %v", err)
				return nil
			}

			_, isLogin, _ := service.GetCurrentCloudBeaverUserID(c)
			if isLogin {
				return respFunc()
			}

			cookies, err := service.Login(user.Name, user.Password)
			if err != nil {
				l.Errorf("login to cloudbeaver failed: %v", err)
				return nil
			}
			for _, cookie := range cookies {
				c.SetCookie(cookie)
			}

			return respFunc()
		}
	}
}

func GraphqlDistributor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().RequestURI != path.Join(service.CbRootUri, service.CbGqlApi) {
				return next(c)
			}
			// copy request body
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			var params *graphql.RawParams
			err := json.Unmarshal(reqBody, &params)
			if err != nil {
				fmt.Println(err)
				return err
			}

			bh, ok := gqlHandlerRouters[params.OperationName]
			if !ok {
				return next(c)
			}

			if bh.disable {
				errMsg := "this feature is prohibited"
				fmt.Printf("%v:%v", errMsg, params.OperationName)
				return c.JSON(http.StatusOK, model.ServerError{
					Message: &errMsg,
				})
			}

			if bh.preprocessing != nil {
				err = bh.preprocessing(c, params)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}

			if bh.useLocalHandler {
				params.ReadTime = graphql.TraceTiming{
					Start: graphql.Now(),
					End:   graphql.Now(),
				}
				ctx := graphql.StartOperationTrace(context.TODO())
				params.Headers = c.Request().Header.Clone()

				var n controller.Next
				var resWrite *responseProcessWriter
				if !bh.needModifyRemoteRes {
					n = func(c echo.Context) ([]byte, error) {
						return nil, next(c)
					}
				} else {
					n = func(c echo.Context) ([]byte, error) {
						resWrite = &responseProcessWriter{tmp: &bytes.Buffer{}, ResponseWriter: c.Response().Writer}
						c.Response().Writer = resWrite
						err := next(c)
						if err != nil {
							return nil, err
						}
						return resWrite.tmp.Bytes(), nil
					}
				}

				g := resolver.NewExecutableSchema(resolver.Config{
					Resolvers: &controller.ResolverImpl{
						Ctx:  c,
						Next: n,
					},
				})

				exec := executor.New(g)

				rc, err := exec.CreateOperationContext(ctx, params)
				if err != nil {
					return err
				}
				responses, ctx := exec.DispatchOperation(ctx, rc)

				res := responses(ctx)
				if res.Errors.Error() != "" {
					return res.Errors
				}
				if !bh.needModifyRemoteRes {
					return nil
				} else {
					header := resWrite.ResponseWriter.Header()
					b, err := json.Marshal(res)
					if err != nil {
						return err
					}
					header.Set("Content-Length", fmt.Sprintf("%d", len(b)))
					_, err = resWrite.ResponseWriter.Write(b)
					return err
				}
			}
			return next(c)
		}
	}
}

type responseProcessWriter struct {
	tmp        *bytes.Buffer
	headerCode int
	http.ResponseWriter
}

func (w *responseProcessWriter) WriteHeader(code int) {
	w.headerCode = code
}

func (w *responseProcessWriter) Write(b []byte) (int, error) {
	return w.tmp.Write(b)
}

func (w *responseProcessWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseProcessWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
