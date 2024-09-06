//go:build enterprise
// +build enterprise

package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type ApiInterfaceInfo struct {
	RouterPath               string
	Method                   string
	OperationType            string
	OperationAction          string
	GetProjectAndContentFunc func(c echo.Context) (projectName, objectName string, err error)
}

var ApiInterfaceInfoList []ApiInterfaceInfo

type ResponseBodyWrite struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseBodyWrite) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseBodyWrite) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.Write([]byte(s))
}

func OperationLogRecord() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			reqIP := c.RealIP()
			path := c.Path()
			newLog := log.NewEntry()
			for _, interfaceInfo := range ApiInterfaceInfoList {
				if c.Request().Method == interfaceInfo.Method && interfaceInfo.RouterPath == path {
					user, err := controller.GetCurrentUser(c, dms.GetUser)
					if err != nil {
						newLog.Errorf("get current error: %s", err)
						return nil
					}
					userName := user.Name

					operationRecord := &model.OperationRecord{
						OperationTime:     time.Now(),
						OperationUserName: userName,
						OperationReqIP:    reqIP,
						OperationTypeName: interfaceInfo.OperationType,
						OperationAction:   interfaceInfo.OperationAction,
					}

					projectName, content, err := interfaceInfo.GetProjectAndContentFunc(c)
					if err != nil {
						newLog.Errorf("get content and project name error: %s", err)
					}

					operationRecord.OperationProjectName = projectName
					operationRecord.OperationContent = content

					respBodyWrite := &ResponseBodyWrite{body: new(bytes.Buffer), ResponseWriter: c.Response().Writer}

					c.Response().Writer = respBodyWrite

					if err = next(c); err != nil {
						c.Error(err)
					}

					resp := respBodyWrite.body.Bytes()
					var respBody map[string]interface{}
					if err := json.Unmarshal(resp, &respBody); err == nil {
						if code, ok := respBody["code"]; ok {
							codeInt := int(code.(float64))
							if codeInt != 0 {
								operationRecord.OperationStatus = model.OperationRecordStatusFailed
							} else {
								operationRecord.OperationStatus = model.OperationRecordStatusSucceeded
							}
						}
					} else {
						operationRecord.OperationStatus = model.OperationRecordStatusFailed
					}

					s := model.GetStorage()
					if err := s.Save(&operationRecord); err != nil {
						newLog.Errorf("save operation record error: %s", err)
						return nil
					}

					return nil
				}
			}

			return next(c)
		}
	}
}
