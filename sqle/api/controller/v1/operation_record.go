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
	return getOperationTypeNameList(c)
}

type GetOperationActionListResV1 struct {
	controller.BaseRes
	Data []OperationActionList `json:"data"`
}

type OperationActionList struct {
	OperationAction string `json:"operation_action"`
	Desc            string `json:"desc"`
}

// GetOperationActionList
// @Summary 获取操作内容列表
// @Description Get operation action list
// @Id getOperationActionList
// @Tags OperationRecord
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetOperationActionListResV1
// @Router /v1/operation_records/operation_actions [get]
func GetOperationActionList(c echo.Context) error {
	return getOperationActionList(c)
}

type GetOperationRecordListReqV1 struct {
	FilterOperateTimeFrom      string  `json:"filter_operate_time_from" query:"filter_operate_time_from"`
	FilterOperateTimeTo        string  `json:"filter_operate_time_to" query:"filter_operate_time_to"`
	FilterOperateProjectName   *string `json:"filter_operate_project_name" query:"filter_operate_project_name"`
	FuzzySearchOperateUserName string  `json:"fuzzy_search_operate_user_name" query:"fuzzy_search_operate_user_name"`
	FilterOperateTypeName      string  `json:"filter_operate_type_name" query:"filter_operate_type_name"`
	FilterOperateAction        string  `json:"filter_operate_action" query:"filter_operate_action"`
	PageIndex                  uint32  `json:"page_index" query:"page_index" valid:"required"`
	PageSize                   uint32  `json:"page_size" query:"page_size" valid:"required"`
}

type GetOperationRecordListResV1 struct {
	controller.BaseRes
	Data      []OperationRecordList `json:"data"`
	TotalNums uint64                `json:"total_nums"`
}

type OperationRecordList struct {
	ID                uint64        `json:"id"`
	OperationTime     *time.Time    `json:"operation_time"`
	OperationUser     OperationUser `json:"operation_user"`
	OperationTypeName string        `json:"operation_type_name"`
	OperationAction   string        `json:"operation_action"`
	OperationContent  string        `json:"operation_content"`
	ProjectName       string        `json:"project_name"`
	Status            string        `json:"status" enums:"succeeded,failed"`
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
// @Param filter_operate_action query string false "filter_operate_action"
// @Param page_index query uint32 true "page_index"
// @Param page_size query uint32 true "page_size"
// @Success 200 {object} v1.GetOperationRecordListResV1
// @Router /v1/operation_records [get]
func GetOperationRecordListV1(c echo.Context) error {
	return getOperationRecordList(c)
}

// GetExportOperationRecordListV1
// @Summary 导出操作记录列表
// @Description Export operation record list
// @Id getExportOperationRecordListV1
// @Tags OperationRecord
// @Security ApiKeyAuth
// @Param filter_operate_time_from query string false "filter_operate_time_from"
// @Param filter_operate_time_to query string false "filter_operate_time_to"
// @Param filter_operate_project_name query string false "filter_operate_project_name"
// @Param fuzzy_search_operate_user_name query string false "fuzzy_search_operate_user_name"
// @Param filter_operate_type_name query string false "filter_operate_type_name"
// @Param filter_operate_action query string false "filter_operate_action"
// @Success 200 {file} file "get export operation record list"
// @Router /v1/operation_records/exports [get]
func GetExportOperationRecordListV1(c echo.Context) error {
	return exportOperationRecordList(c)
}
