package accesstoken

import (
	"context"
	"fmt"
	"net/http"

	jwtPkg "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/labstack/echo/v4"
)

const AccessTokenLogin = "access_token_login"

func CheckLatestAccessToken(dmsAddress string, getTokenDetail func(c jwtPkg.EchoContextGetter) (*jwtPkg.TokenDetail, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenDetail, err := getTokenDetail(c)

			if err != nil {
				echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("get token detail failed, err:%v", err))
				return err
			}

			if tokenDetail.TokenStr == "" {
				return next(c)
			}

			// LoginType为空，不需要校验access token
			if tokenDetail.LoginType == "" {
				return next(c)
			}

			if tokenDetail.LoginType != AccessTokenLogin {
				return echo.NewHTTPError(http.StatusUnauthorized, "access token login type is error")
			}

			userInfo, err := dmsobject.GetUser(context.TODO(), tokenDetail.UID, dmsAddress)
			if err != nil {
				return err
			}
			if userInfo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "access token: cannot get user info")
			}

			if userInfo.AccessTokenInfo.AccessToken != tokenDetail.TokenStr {
				return echo.NewHTTPError(http.StatusUnauthorized, "access token is not latest")
			}
			return next(c)
		}
	}
}
