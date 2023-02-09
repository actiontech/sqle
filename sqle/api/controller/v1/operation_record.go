package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetOperationTypeNamesListResV1 struct {
	controller.BaseRes
	Data []OperationTypeNameList `json:"data"`
}

type OperationTypeNameList struct {
	OperationTypeName string `json:"operation_type_name"`
	Desc              string `json:"desc"`
}

// GetOperationTypeNameList
// @Summary 获取操作类型名列表
// @Description Get operation type name list
// @Id GetOperationTypeNameList
// @Tags OperationRecord
// @Security ApiKeyAuth
// @Success 200 {object} GetOperationTypeNamesListResV1
// @Router /v1/operation_records/operation_type_names [get]
func GetOperationTypeNameList(c echo.Context) error {
	return nil
}

type GetOperationContentListResV1 struct {
	controller.BaseRes
	Data []OperationContentList `json:"data"`
}

type OperationContentList struct {
	OperationContent string `json:"operation_content"`
	Desc             string `json:"desc"`
}

// GetOperationContentList
// @Summary 获取操作内容列表
// @Description Get operation content list
// @Id getOperationContentList
// @Tags OperationRecord
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetOperationContentListResV1
// @Router /v1/operation_records/operation_contents [get]
func GetOperationContentList(c echo.Context) error {
	return nil
}

type GetOperationRecordListResV1 struct {
	controller.BaseRes
	Data      []OperationRecordList `json:"data"`
	TotalNums uint64                `json:"total_nums"`
}

type OperationRecordList struct {
	ID                  uint64        `json:"id"`
	OperationTime       *time.Time    `json:"operation_time"`
	OperationUser       OperationUser `json:"operation_user"`
	OperationTypeName   string        `json:"operation_type_name"`
	OperationContent    string        `json:"operation_content"`
	OperationObjectName string        `json:"operation_object_name"`
	ProjectName         string        `json:"project_name"`
	Status              string        `json:"status" enums:"success,fail"`
}

type OperationUser struct {
	UserName string `json:"user_name"`
	IP       string `json:"ip"`
}

// GetOperationRecordListV1
// @Summary 获取操作记录列表
// @Description Get operation record list
// @Id getOperationRecordListV1
// @Tags OperationRecord
// @Security ApiKeyAuth
// @Param filter_operate_time_from query string false "filter_operate_time_from"
// @Param filter_operate_time_to query string false "filter_operate_time_to"
// @Param filter_operate_project_name query string false "filter_operate_project_name"
// @Param fuzzy_search_operate_user_name query string false "fuzzy_search_operate_user_name"
// @Param filter_operate_type_name query string false "filter_operate_type_name"
// @Param filter_operate_content query string false "filter_operate_content"
// @Param page_index query uint32 true "page_index"
// @Param page_size query uint32 true "page_size"
// @Success 200 {object} v1.GetOperationRecordListResV1
// @Router /v1/operation_records [get]
func GetOperationRecordListV1(c echo.Context) error {
	return nil
}
