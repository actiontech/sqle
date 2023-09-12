package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type CreateSQLAuditRecordReqV1 struct {
	DbType         string `json:"db_type" form:"db_type" example:"MySQL"`
	InstanceName   string `json:"instance_name" form:"instance_name" example:"inst_1"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema" example:"db1"`
	Sqls           string `json:"sqls" form:"sqls" example:"alter table tb1 drop columns c1; select * from tb"`
}

type CreateSQLAuditRecordResV1 struct {
	controller.BaseRes
	Data *SQLAuditRecordResData `json:"data"`
}

type SQLAuditRecordResData struct {
	Id   string          `json:"sql_audit_record_id"`
	Task *AuditTaskResV1 `json:"task"`
}

// CreateSQLAuditRecord
// @Summary SQL审核
// @Id CreateSQLAuditRecordV1
// @Description SQL audit
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is ZIP file that sql will be parsed from xml or sql file inside it.
// @Accept mpfd
// @Produce json
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name formData string false "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param db_type formData string false "db type of instance"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.CreateSQLAuditRecordResV1
// @router /v1/projects/{project_name}/sql_audit_record [post]
func CreateSQLAuditRecord(c echo.Context) error {
	return nil
func buildOnlineTaskForAudit(c echo.Context, s *model.Storage, userId uint, instanceName, instanceSchema, projectName, sourceType, sqls string) (*model.Task, error) {
	instance, exist, err := s.GetInstanceByNameAndProjectName(instanceName, projectName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrInstanceNoAccess
	}
	can, err := checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return nil, err
	}
	if !can {
		return nil, ErrInstanceNoAccess
	}

	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	if err := plugin.Ping(context.TODO()); err != nil {
		return nil, err
	}

	task := &model.Task{
		Schema:       instanceSchema,
		InstanceId:   instance.ID,
		Instance:     instance,
		CreateUserId: userId,
		ExecuteSQLs:  []*model.ExecuteSQL{},
		SQLSource:    sourceType,
		DBType:       instance.DbType,
	}
	createAt := time.Now()
	task.CreatedAt = createAt

	nodes, err := plugin.Parse(context.TODO(), sqls)
	if err != nil {
		return nil, err
	}
	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text,
			},
		})
	}
	return task, nil
}
}

type UpdateSQLAuditRecordReqV1 struct {
	SQLAuditRecordId string    `json:"sql_audit_record_id" valid:"required"`
	Tags             *[]string `json:"tags"`
}

// UpdateSQLAuditRecordV1
// @Summary 更新SQL审核记录
// @Description update SQL audit record
// @Accept json
// @Id updateSQLAuditRecordV1
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param param body v1.UpdateSQLAuditRecordReqV1 true "update SQL audit record"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_audit_record/{sql_audit_record_id} [patch]
func UpdateSQLAuditRecordV1(c echo.Context) error {
	return nil
}

type GetSQLAuditRecordsReqV1 struct {
	FuzzySearchSQLAuditRecordId string `json:"fuzzy_search_sql_audit_record_id" query:"fuzzy_search_sql_audit_record_id"`
	FilterSQLAuditStatus        string `json:"filter_sql_audit_status" query:"filter_sql_audit_status" enums:"auditing,successfully,"`
	FilterTags                  string `json:"filter_tags" query:"filter_tags"`
}

type SQLAuditRecord struct {
	Creator          string                 `json:"creator"`
	SQLAuditRecordId uint                   `json:"sql_audit_record_id"`
	SQLAuditStatus   string                 `json:"sql_audit_status"`
	Tags             []string               `json:"tags"`
	Instance         SQLAuditRecordInstance `json:"instance"`
	Task             AuditTaskResV1         `json:"task"`
}

type SQLAuditRecordInstance struct {
	Host string `json:"db_host" example:"10.10.10.10"`
	Port string `json:"db_port" example:"3306"`
}

type GetSQLAuditRecordsResV1 struct {
	controller.BaseRes
	Data      []SQLAuditRecord `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

// GetSQLAuditRecordsV1
// @Summary 获取SQL审核记录列表
// @Description get sql audit records
// @Tags sql_audit_record
// @Id getSQLAuditRecordsV1
// @Security ApiKeyAuth
// @Param fuzzy_search_tags query string false "fuzzy search tags"
// @Param filter_sql_audit_status query string false "filter sql audit status" Enums(auditing,successfully)
// @Param filter_instance_name query string false "filter instance name"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetSQLAuditRecordsResV1
// @router /v1/projects/{project_name}/sql_audit_record [get]
func GetSQLAuditRecordsV1(c echo.Context) error {
	return nil
}

type GetSQLAuditRecordResV1 struct {
	controller.BaseRes
	Data SQLAuditRecord `json:"data"`
}

// GetSQLAuditRecordV1
// @Summary 获取SQL审核记录信息
// @Description get sql audit record info
// @Tags sql_audit_record
// @Id getSQLAuditRecordV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_audit_record_id path string true "sql audit record id"
// @Success 200 {object} v1.GetSQLAuditRecordResV1
// @router /v1/projects/{project_name}/sql_audit_record/{sql_audit_record_id} [get]
func GetSQLAuditRecordV1(c echo.Context) error {
	return nil
}

type GetSQLAuditRecordTagTipsResV1 struct {
	controller.BaseRes
	Tags []string `json:"data"`
}

// GetSQLAuditRecordTagTipsV1
// @Summary 获取SQL审核记录标签列表
// @Description get sql audit record tag tips
// @Tags sql_audit_record
// @Id GetSQLAuditRecordTagTipsV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetSQLAuditRecordTagTipsResV1
// @router /v1/projects/{project_name}/sql_audit_record/tag_tips [get]
func GetSQLAuditRecordTagTipsV1(c echo.Context) error {
	return nil
}
