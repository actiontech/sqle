package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

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
		return err
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

func GetCurrentUser(c echo.Context) (*model.User, error) {
	key := "current_user"
	currentUser := c.Get(key)
	if currentUser != nil {
		if user, ok := currentUser.(*model.User); ok {
			return user, nil
		}
	}
	s := model.GetStorage()
	user, exist, err := s.GetUserByName(GetUserName(c))
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(errors.DataNotExist,
			fmt.Errorf("current user is not exist"))
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

// ReadFileContent read content from http body by name if file exist,
// the name is a http form data key, not file name.
func ReadFileContent(c echo.Context, name string) (content string, fileExist bool, err error) {
	file, err := c.FormFile(name)
	if err == http.ErrMissingFile {
		return "", false, nil
	}
	if err != nil {
		return "", false, errors.New(errors.ReadUploadFileError, err)
	}
	src, err := file.Open()
	if err != nil {
		return "", false, errors.New(errors.ReadUploadFileError, err)
	}
	defer src.Close()
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return "", false, errors.New(errors.ReadUploadFileError, err)
	}
	return string(data), true, nil
}

func IsUserCanBeDisabled(editorUserName, editedUserName string) (err error) {

	if editorUserName == editedUserName {
		return errors.NewDataInvalidErr("user<%v> can not disable or enable self", editorUserName)
	}

	// admin user can not be disabled.
	if model.IsDefaultAdminUser(editedUserName) {
		return errors.NewDataInvalidErr("admin user can not be disabled or enabled")
	}

	return nil
}
