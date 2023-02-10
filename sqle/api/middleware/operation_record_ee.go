//go:build enterprise
// +build enterprise

package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

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
			reqIP := c.Request().Host
			path := c.Path()
			newLog := log.NewEntry()
			for _, interfaceInfo := range v1.ApiInterfaceInfoList {
				if c.Request().Method == interfaceInfo.Method && interfaceInfo.RouterPath == path {
					userName := controller.GetUserName(c)

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
								operationRecord.OperationStatus = model.OperationRecordStatusFail
							} else {
								operationRecord.OperationStatus = model.OperationRecordStatusSuccess
							}
						}
					} else {
						operationRecord.OperationStatus = model.OperationRecordStatusFail
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
