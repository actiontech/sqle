package cloudbeaver_wrapper

// import (
// 	"bufio"
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io/ioutil"
// 	"net"
// 	"net/http"
// 	"net/url"
// 	"path"
// 	"strconv"
// 	"sync"

// 	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/controller"
// 	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
// 	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/resolver"
// 	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
// 	"github.com/actiontech/sqle/sqle/log"
// 	sqleModel "github.com/actiontech/sqle/sqle/model"
// 	"github.com/actiontech/sqle/sqle/utils"

// 	"github.com/99designs/gqlgen/graphql"
// 	"github.com/99designs/gqlgen/graphql/executor"
// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// type gqlBehavior struct {
// 	useLocalHandler     bool
// 	needModifyRemoteRes bool
// 	disable             bool
// 	// 预处理主要用于在真正使用前处理前端传递的参数, 比如需要接收int, 但收到float, 则可以在此处调整参数类型
// 	preprocessing func(ctx echo.Context, params *graphql.RawParams) error
// }

// var gqlHandlerRouters = map[string] /* gql operation name */ gqlBehavior{
// 	"asyncSqlExecuteQuery": {
// 		useLocalHandler:     true,
// 		needModifyRemoteRes: false,
// 		preprocessing: func(ctx echo.Context, params *graphql.RawParams) (err error) {
// 			// json中没有int类型, 这将导致执行json.Unmarshal()时int会被当作float64, 从而导致后面出现类型错误的异常
// 			if filter, ok := params.Variables["filter"].(map[string]interface{}); ok {
// 				if filter["limit"] != nil {
// 					params.Variables["filter"].(map[string]interface{})["limit"], err = strconv.Atoi(fmt.Sprintf("%v", params.Variables["filter"].(map[string]interface{})["limit"]))
// 				}
// 			}
// 			return err
// 		},
// 	},
// 	"getActiveUser": {
// 		useLocalHandler:     true,
// 		needModifyRemoteRes: true,
// 	}, "authLogout": {
// 		disable: true,
// 	}, "authLogin": {
// 		disable: true,
// 	}, "configureServer": {
// 		disable: true,
// 	}, "createUser": {
// 		disable: true,
// 	}, "setUserCredentials": {
// 		disable: true,
// 	}, "enableUser": {
// 		disable: true,
// 	}, "grantUserRole": {
// 		disable: true,
// 	}, "setConnections": {
// 		disable: true,
// 	}, "saveUserMetaParameters": {
// 		disable: true,
// 	}, "deleteUser": {
// 		disable: true,
// 	}, "createRole": {
// 		disable: true,
// 	}, "updateRole": {
// 		disable: true,
// 	}, "deleteRole": {
// 		disable: true,
// 	}, "authChangeLocalPassword": {
// 		disable: true,
// 	},
// }

// func StartApp(e *echo.Echo) error {
// 	if !service.IsCloudBeaverConfigured() {
// 		return nil
// 	}
// 	cfg := service.GetSQLQueryConfig()
// 	protocol := "http"
// 	if cfg.EnableHttps {
// 		protocol = "https"
// 	}
// 	url2, err := url.Parse(fmt.Sprintf("%v://%v:%v", protocol, cfg.CloudBeaverHost, cfg.CloudBeaverPort))
// 	if err != nil {
// 		return err
// 	}
// 	targets := []*middleware.ProxyTarget{
// 		{
// 			URL: url2,
// 		},
// 	}

// 	err = service.InitGQLVersion()
// 	if err != nil {
// 		return err
// 	}

// 	q := e.Group(service.CbRootUri)

// 	q.Use(TriggerLogin())
// 	q.Use(GraphqlDistributor())
// 	q.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
// 		Skipper:  middleware.DefaultSkipper,
// 		Balancer: middleware.NewRandomBalancer(targets),
// 	}))

// 	return nil
// }

// var (
// 	sqleTokenToCBSessionId = make(map[string]string)
// 	tokenMapMutex          = &sync.Mutex{}
// )

// func getCBSessionIdBySqleToken(token string) string {
// 	tokenMapMutex.Lock()
// 	defer tokenMapMutex.Unlock()
// 	return sqleTokenToCBSessionId[token]
// }

// func setCBSessionIdBySqleToken(token, cbSessionId string) {
// 	tokenMapMutex.Lock()
// 	defer tokenMapMutex.Unlock()
// 	sqleTokenToCBSessionId[token] = cbSessionId
// }

// func UnbindCBSessionIdBySqleToken(token string) {
// 	tokenMapMutex.Lock()
// 	defer tokenMapMutex.Unlock()
// 	delete(sqleTokenToCBSessionId, token)
// }

// // 如果当前用户没有登录cloudbeaver，则登录
// func TriggerLogin() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			var sqleToken string
// 			// 根据cookie中的sqle-token查找对应用户的cb-session-id
// 			for _, c := range c.Cookies() {
// 				if c.Name == "sqle-token" {
// 					sqleToken = c.Value
// 					break
// 				}
// 			}
// 			if sqleToken == "" {
// 				// 没有找到sqle-token，有可能是用户直接通过url访问cb页面，但没有登录sqle
// 				return c.Redirect(http.StatusFound, "/login?target=/sqlQuery")
// 			}
// 			CBSessionId := getCBSessionIdBySqleToken(sqleToken)
// 			if CBSessionId != "" {
// 				// todo 处理sessionId超时的情况
// 				c.Request().Header.Set("Cookie", "cb-session-id="+CBSessionId)
// 				return next(c)
// 			}

// 			// CBSessionId不存在认为当前用户没有登录cb，登录cb
// 			l := log.NewEntry().WithField("action", "trigger cloudbeaver login")
// 			userName, err := utils.GetUserNameFromJWTToken(sqleToken)
// 			if err != nil {
// 				l.Errorf("get user name from token failed: %v", err)
// 				return errors.New("get user name to login failed")
// 			}
// 			s := sqleModel.GetStorage()
// 			user, _, err := s.GetUserByName(userName)
// 			if err != nil {
// 				l.Errorf("get user info err: %v", err)
// 				return err
// 			}

// 			cbUser := service.GenerateCloudBeaverUserName(userName)
// 			// 同步信息
// 			if err = service.SyncCurrentUser(cbUser); err != nil {
// 				l.Errorf("sync cloudbeaver user %v info failed: %v", cbUser, err)
// 			}
// 			err = service.SyncUserBindInstance(cbUser)
// 			if err != nil {
// 				l.Errorf("sync cloudbeaver user %v bind instance failed: %v", cbUser, err)
// 			}
// 			cookies, err := service.LoginToCBServer(cbUser, user.Password)
// 			if err != nil {
// 				l.Errorf("login to cloudbeaver failed: %v", err)
// 				return err
// 			}

// 			// 添加sqle和cb的用户映射
// 			for _, ck := range cookies {
// 				if ck.Name == "cb-session-id" {
// 					setCBSessionIdBySqleToken(sqleToken, ck.Value)
// 				}
// 			}

// 			return next(c)
// 		}
// 	}
// }

// func GraphqlDistributor() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			if c.Request().RequestURI != path.Join(service.CbRootUri, service.CbGqlApi) {
// 				return next(c)
// 			}
// 			// copy request body
// 			reqBody := []byte{}
// 			if c.Request().Body != nil { // Read
// 				reqBody, _ = ioutil.ReadAll(c.Request().Body)
// 			}
// 			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

// 			var params *graphql.RawParams
// 			err := json.Unmarshal(reqBody, &params)
// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}

// 			bh, ok := gqlHandlerRouters[params.OperationName]
// 			if !ok {
// 				return next(c)
// 			}

// 			if bh.disable {
// 				errMsg := "this feature is prohibited"
// 				fmt.Printf("%v:%v", errMsg, params.OperationName)
// 				return c.JSON(http.StatusOK, model.ServerError{
// 					Message: &errMsg,
// 				})
// 			}

// 			if bh.preprocessing != nil {
// 				err = bh.preprocessing(c, params)
// 				if err != nil {
// 					fmt.Println(err)
// 					return err
// 				}
// 			}

// 			if bh.useLocalHandler {
// 				params.ReadTime = graphql.TraceTiming{
// 					Start: graphql.Now(),
// 					End:   graphql.Now(),
// 				}
// 				ctx := graphql.StartOperationTrace(context.TODO())
// 				params.Headers = c.Request().Header.Clone()

// 				var n controller.Next
// 				var resWrite *responseProcessWriter
// 				if !bh.needModifyRemoteRes {
// 					n = func(c echo.Context) ([]byte, error) {
// 						return nil, next(c)
// 					}
// 				} else {
// 					n = func(c echo.Context) ([]byte, error) {
// 						resWrite = &responseProcessWriter{tmp: &bytes.Buffer{}, ResponseWriter: c.Response().Writer}
// 						c.Response().Writer = resWrite
// 						err := next(c)
// 						if err != nil {
// 							return nil, err
// 						}
// 						return resWrite.tmp.Bytes(), nil
// 					}
// 				}

// 				g := resolver.NewExecutableSchema(resolver.Config{
// 					Resolvers: &controller.ResolverImpl{
// 						Ctx:  c,
// 						Next: n,
// 					},
// 				})

// 				exec := executor.New(g)

// 				rc, err := exec.CreateOperationContext(ctx, params)
// 				if err != nil {
// 					return err
// 				}
// 				responses, ctx := exec.DispatchOperation(ctx, rc)

// 				res := responses(ctx)
// 				if res.Errors.Error() != "" {
// 					return res.Errors
// 				}
// 				if !bh.needModifyRemoteRes {
// 					return nil
// 				} else {
// 					header := resWrite.ResponseWriter.Header()
// 					b, err := json.Marshal(res)
// 					if err != nil {
// 						return err
// 					}
// 					header.Set("Content-Length", fmt.Sprintf("%d", len(b)))
// 					_, err = resWrite.ResponseWriter.Write(b)
// 					return err
// 				}
// 			}
// 			return next(c)
// 		}
// 	}
// }

// type responseProcessWriter struct {
// 	tmp        *bytes.Buffer
// 	headerCode int
// 	http.ResponseWriter
// }

// func (w *responseProcessWriter) WriteHeader(code int) {
// 	w.headerCode = code
// }

// func (w *responseProcessWriter) Write(b []byte) (int, error) {
// 	return w.tmp.Write(b)
// }

// func (w *responseProcessWriter) Flush() {
// 	w.ResponseWriter.(http.Flusher).Flush()
// }

// func (w *responseProcessWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// 	return w.ResponseWriter.(http.Hijacker).Hijack()
// }
