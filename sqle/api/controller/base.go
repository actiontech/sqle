package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	dmsJWT "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var dmsServerAddress string

func GetDMSServerAddress() string {
	return dmsServerAddress
}
func InitDMSServerAddress(addr string) {
	dmsServerAddress = addr
}

type BaseRes struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message" example:"ok"`
}

func NewBaseReq(err error) BaseRes {
	res := BaseRes{}
	switch e := err.(type) {
	case *errors.CodeError:
		res.Code = e.Code()
		res.Message = e.Error()
	default:
		if err == nil {
			res.Code = 0
			res.Message = "ok"
		} else {
			res.Code = -1
			res.Message = e.Error()
		}
	}
	return res
}

func BindAndValidateReq(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		err := fmt.Errorf("bind request error: %v, please check request body", err)
		return errors.HttpRequestFormatErrWrapper(err)
	}

	if err := Validate(i); err != nil {
		return errors.New(errors.DataInvalid, err)
	}
	return nil
}

func GetUserName(c echo.Context) string {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return ""
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	return claims["name"].(string)
}

func GetUserID(c echo.Context) string {
	uidStr, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
	if err != nil {
		return ""
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", uid)
}

func GetCurrentUser(c echo.Context, getUser func(context.Context, string, string) (*model.User, error)) (*model.User, error) {
	key := "current_user"
	currentUser := c.Get(key)
	if currentUser != nil {
		if user, ok := currentUser.(*model.User); ok {
			return user, nil
		}
	}
	uidStr := GetUserID(c)
	user, err := getUser(c.Request().Context(), uidStr, GetDMSServerAddress())
	if err != nil {
		return nil, err
	}
	c.Set(key, user)
	return user, nil
}

func JSONBaseErrorReq(c echo.Context, err error) error {
	return c.JSON(http.StatusOK, NewBaseReq(err))
}

func JSONNewNotImplementedErr(c echo.Context) error {
	return c.JSON(http.StatusOK,
		NewBaseReq(errors.NewNotImplementedError("not implemented yet")))
}

func JSONNewDataExistErr(c echo.Context, format string, a ...interface{}) error {
	return c.JSON(http.StatusOK,
		NewBaseReq(errors.New(errors.DataExist, fmt.Errorf(format, a...))))
}

func JSONNewDataNotExistErr(c echo.Context, format string, a ...interface{}) error {
	return c.JSON(http.StatusOK,
		NewBaseReq(errors.New(errors.DataNotExist, fmt.Errorf(format, a...))))
}

func JSONOnlySupportForEnterpriseVersionErr(c echo.Context) error {
	return c.JSON(http.StatusOK, NewBaseReq(errors.NewOnlySupportForEnterpriseVersion()))
}

// ReadFile read content from http body by name if file exist,
// the name is a http form data key, not file name.
func ReadFile(c echo.Context, name string) (fileName, content string, fileExist bool, err error) {
	file, err := c.FormFile(name)
	if err == http.ErrMissingFile {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, errors.New(errors.ReadUploadFileError, err)
	}
	src, err := file.Open()
	if err != nil {
		return "", "", false, errors.New(errors.ReadUploadFileError, err)
	}
	defer src.Close()
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return "", "", false, errors.New(errors.ReadUploadFileError, err)
	}
	return file.Filename, string(data), true, nil
}

// subjectUser should be admin user.
func CanThisUserBeDisabled(subjectUser, objectUser string) (err error) {

	if subjectUser == objectUser {
		return errors.NewDataInvalidErr("user<%v> can not disable or enable self", subjectUser)
	}

	// admin user can not be disabled.
	if model.IsDefaultAdminUser(objectUser) {
		return errors.NewDataInvalidErr("admin user can not be disabled or enabled")
	}

	return nil
}

func GetLimitAndOffset(pageIndex, pageSize uint32) (limit, offset uint32) {
	if pageIndex >= 1 {
		offset = (pageIndex - 1) * pageSize
	}
	return pageSize, offset
}
