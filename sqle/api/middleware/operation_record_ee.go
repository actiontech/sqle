//go:build enterprise
// +build enterprise

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type apiInterfaceInfo struct {
	reg                     *regexp.Regexp
	method                  string
	operationType           string
	operationContent        string
	getProjectAndObjectFunc func(c echo.Context) (projectName, objectName string, err error)
}

var apiInterfaceInfoList = []apiInterfaceInfo{
	{
		reg:                     regexp.MustCompile("/v1/projects"),
		method:                  http.MethodPost,
		operationType:           model.OperationRecordProjectManageType,
		operationContent:        model.OperationRecordCreateProjectContent,
		getProjectAndObjectFunc: getProjectAndObjectFromCreateProject,
	},
}

func getProjectAndObjectFromCreateProject(c echo.Context) (string, string, error) {
	req := new(v1.CreateProjectReqV1)

	reqBody, err := getReqBodyBytes(c)
	if err != nil {
		return "", "", err
	}

	if err := json.Unmarshal(reqBody, req); err != nil {
		return "", "", err
	}

	if err := controller.Validate(req); err != nil {
		return "", "", err
	}

	return model.OperationRecordPlatform, req.Name, nil
}

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
			path := c.Request().URL.Path
			newLog := log.NewEntry()
			for _, interfaceInfo := range apiInterfaceInfoList {
				if c.Request().Method == interfaceInfo.method && interfaceInfo.reg.MatchString(path) {
					userName := controller.GetUserName(c)

					operationRecord := &model.OperationRecord{
						OperationTime: time.Now(),
						UserName:      userName,
						IP:            reqIP,
						TypeName:      interfaceInfo.operationType,
						Content:       interfaceInfo.operationContent,
					}

					projectName, objectName, err := interfaceInfo.getProjectAndObjectFunc(c)
					if err != nil {
						newLog.Errorf("get object and project name error: %s", err)
					}

					operationRecord.ProjectName = projectName
					operationRecord.ObjectName = objectName

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
								operationRecord.Status = model.OperationRecordFailStatus
							} else {
								operationRecord.Status = model.OperationRecordSuccessStatus
							}
						}
					} else {
						operationRecord.Status = model.OperationRecordFailStatus
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

func getReqBodyBytes(c echo.Context) ([]byte, error) {
	var bodyBytes []byte
	var err error

	if c.Request().Body != nil {
		bodyBytes, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return nil, err
		}

		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return bodyBytes, nil
	}

	return nil, fmt.Errorf("request body is nil")
}
