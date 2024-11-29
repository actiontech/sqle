package v1

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	e "errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	mybatis_parser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var ErrTooManyDataSource = errors.New(errors.DataConflict, fmt.Errorf("the number of data sources must be less than %v", MaximumDataSourceNum))

type CreateAuditTaskReqV1 struct {
	InstanceName    string `json:"instance_name" form:"instance_name" example:"inst_1" valid:"required"`
	InstanceSchema  string `json:"instance_schema" form:"instance_schema" example:"db1"`
	Sql             string `json:"sql" form:"sql" example:"alter table tb1 drop columns c1"`
	ExecMode        string `json:"exec_mode" form:"exec_mode" enums:"sql_file,sqls"`
	EnableBackup    bool   `json:"enable_backup" form:"enable_backup"`
	FileOrderMethod string `json:"file_order_method" form:"file_order_method"`
}

type GetAuditTaskResV1 struct {
	controller.BaseRes
	Data *AuditTaskResV1 `json:"data"`
}

type AuditTaskResV1 struct {
	Id                         uint            `json:"task_id"`
	InstanceName               string          `json:"instance_name"`
	InstanceDbType             string          `json:"instance_db_type"`
	InstanceSchema             string          `json:"instance_schema" example:"db1"`
	AuditLevel                 string          `json:"audit_level" enums:"normal,notice,warn,error,"`
	Score                      int32           `json:"score"`
	PassRate                   float64         `json:"pass_rate"`
	Status                     string          `json:"status" enums:"initialized,audited,executing,exec_success,exec_failed,manually_executed"`
	SQLSource                  string          `json:"sql_source" enums:"form_data,sql_file,mybatis_xml_file,audit_plan,zip_file,git_repository"`
	ExecStartTime              *time.Time      `json:"exec_start_time,omitempty"`
	ExecEndTime                *time.Time      `json:"exec_end_time,omitempty"`
	FileOrderMethod            string          `json:"file_order_method,omitempty"`
	ExecMode                   string          `json:"exec_mode,omitempty"`
	EnableBackup               bool            `json:"enable_backup"`
	BackupConflictWithInstance bool            `json:"backup_conflict_with_instance"` // 当数据源备份开启，工单备份关闭，则需要提示审核人工单备份策略与数据源备份策略不一致
	AuditFiles                 []AuditFileResp `json:"audit_files,omitempty"`
}

type AuditFileResp struct {
	FileName string `json:"file_name"`
}

func convertTaskToRes(task *model.Task) *AuditTaskResV1 {
	return &AuditTaskResV1{
		Id:                         task.ID,
		InstanceName:               task.InstanceName(),
		InstanceDbType:             task.DBType,
		InstanceSchema:             task.Schema,
		AuditLevel:                 task.AuditLevel,
		Score:                      task.Score,
		PassRate:                   task.PassRate,
		Status:                     task.Status,
		SQLSource:                  task.SQLSource,
		ExecStartTime:              task.ExecStartAt,
		ExecEndTime:                task.ExecEndAt,
		ExecMode:                   task.ExecMode,
		EnableBackup:               task.EnableBackup,
		BackupConflictWithInstance: server.BackupService{}.IsBackupConflictWithInstance(task.EnableBackup, task.InstanceEnableBackup),
		FileOrderMethod:            task.FileOrderMethod,
		AuditFiles:                 convertToAuditFileResp(task.AuditFiles),
	}
}
func convertToAuditFileResp(files []*model.AuditFile) []AuditFileResp {
	fileResp := make([]AuditFileResp, 0, len(files))
	for _, file := range files {
		fileResp = append(fileResp, AuditFileResp{
			FileName: file.FileName,
		})
	}
	return fileResp
}

const (
	InputSQLFileName        = "input_sql_file"
	InputMyBatisXMLFileName = "input_mybatis_xml_file"
	InputZipFileName        = "input_zip_file"
	InputFileFromGit        = "input_file_from_git"
	GitHttpURL              = "git_http_url"
	GitUserName             = "git_user_name"
	GitPassword             = "git_user_password"
	ZIPFileExtension        = ".zip"
)

func getSQLFromFile(c echo.Context) (getSQLFromFileResp, error) {
	// Read it from sql file.
	fileName, sqlsFromSQLFile, exist, err := controller.ReadFile(c, InputSQLFileName)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		return getSQLFromFileResp{
			SourceType: model.TaskSQLSourceFromSQLFile,
			SQLsFromSQLFiles: []SQLsFromSQLFile{{
				FilePath: fileName,
				SQLs:     sqlsFromSQLFile}},
		}, nil
	}

	// If sql_file is not exist, read it from mybatis xml file.
	fileName, data, exist, err := controller.ReadFile(c, InputMyBatisXMLFileName)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		sqls, err := mybatis_parser.ParseXMLs([]mybatis_parser.XmlFile{{Content: data}}, mybatis_parser.SkipErrorQuery, mybatis_parser.RestoreOriginSql)
		if err != nil {
			return getSQLFromFileResp{}, errors.New(errors.ParseMyBatisXMLFileError, err)
		}
		sqlsFromXMLs := make([]SQLFromXML, len(sqls))
		for i := range sqls {
			sqlsFromXMLs[i] = SQLFromXML{
				FilePath:  fileName,
				StartLine: sqls[i].StartLine,
				SQL:       sqls[i].SQL,
			}
		}
		return getSQLFromFileResp{
			SourceType:   model.TaskSQLSourceFromMyBatisXMLFile,
			SQLsFromXMLs: sqlsFromXMLs,
		}, nil
	}

	// If mybatis xml file is not exist, read it from zip file.
	sqlsFromSQLFiles, sqlsFromXML, exist, err := getSqlsFromZip(c)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		return getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromZipFile,
			SQLsFromSQLFiles: sqlsFromSQLFiles,
			SQLsFromXMLs:     sqlsFromXML,
		}, nil
	}

	// If zip file is not exist, read it from git repository
	sqlsFromSQLFiles, sqlsFromJavaFiles, sqlsFromXMLs, exist, err := getSqlsFromGit(c)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		return getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromGitRepository,
			SQLsFromSQLFiles: append(sqlsFromSQLFiles, sqlsFromJavaFiles...),
			SQLsFromXMLs:     sqlsFromXMLs,
		}, nil
	}
	return getSQLFromFileResp{}, errors.New(errors.DataInvalid, fmt.Errorf("input sql is empty"))
}

func saveFileFromContext(c echo.Context) ([]*model.AuditFile, error) {
	fileHeader, fileType, err := getFileHeaderFromContext(c)
	if err != nil {
		return nil, err
	}
	if !isSupportFileType(fileType) {
		return nil, nil
	}
	multipartFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer multipartFile.Close()

	err = utils.EnsureFilePathWithPermission(model.FixFilePath, utils.OwnerPrivilegedAccessMode)
	if err != nil {
		return nil, err
	}
	uniqueName := model.GenUniqueFileName()
	err = utils.SaveFile(multipartFile, model.DefaultFilePath(uniqueName))
	if err != nil {
		return nil, err
	}
	auditFiles := []*model.AuditFile{
		model.NewFileRecord(0, 1, fileHeader.Filename, uniqueName),
	}
	if strings.HasSuffix(fileHeader.Filename, ".zip") {
		auditFiles[0].ExecOrder = 0
		auditFilesInZip, err := getFileRecordsFromZip(multipartFile, fileHeader)
		if err != nil {
			return nil, err
		}
		auditFiles = append(auditFiles, auditFilesInZip...)
	}
	return auditFiles, nil
}

func getFileRecordsFromZip(multipartFile multipart.File, fileHeader *multipart.FileHeader) ([]*model.AuditFile, error) {
	r, err := zip.NewReader(multipartFile, fileHeader.Size)
	if err != nil {
		return nil, err
	}
	var auditFiles []*model.AuditFile
	var execOrder uint = 1
	for _, srcFile := range r.File {
		// skip empty file and folder
		if srcFile == nil || srcFile.FileInfo().IsDir() {
			continue
		}
		fullName := srcFile.FileHeader.Name // full name with relative path to zip file
		if srcFile.NonUTF8 {
			utf8NameByte, err := utils.ConvertToUtf8([]byte(fullName))
			if err != nil {
				if e.Is(err, utils.ErrUnknownEncoding) {
					return nil, e.New("the file name contains unrecognized characters. Please ensure the file name is encoded in UTF-8 or use an English file name")
				}
				return nil, err
			} else {
				fullName = string(utf8NameByte)
			}
		}
		if strings.HasSuffix(fullName, ".sql") {
			auditFiles = append(auditFiles, model.NewFileRecord(0, execOrder, fullName, model.GenUniqueFileName()))
			execOrder++
		}
	}
	return auditFiles, nil
}

func isSupportFileType(fileType string) bool {
	return fileType == InputSQLFileName || fileType == InputZipFileName || fileType == InputMyBatisXMLFileName
}

func getFileHeaderFromContext(c echo.Context) (fileHeader *multipart.FileHeader, fileType string, err error) {
	if c.FormValue(GitHttpURL) != "" {
		return nil, InputFileFromGit, nil
	}
	fileTypes := []string{
		InputSQLFileName,
		InputMyBatisXMLFileName,
		InputZipFileName,
	}
	for _, fileType = range fileTypes {
		fileHeader, err = c.FormFile(fileType)
		if err == http.ErrMissingFile {
			continue
		}
		if err != nil {
			return nil, fileType, errors.New(errors.ReadUploadFileError, err)
		}
		if fileHeader != nil {
			return fileHeader, fileType, nil
		}
	}
	return nil, "", fmt.Errorf("unknown input file type")
}

// @Summary 创建Sql扫描任务并提交审核
// @Description create and audit a task, you can upload sql content in three ways, any one can be used, but only one is effective.
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Accept mpfd
// @Produce json
// @Tags task
// @Id createAndAuditTaskV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name formData string true "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param enable_backup formData bool false "enable backup"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Param exec_mode formData string false "exec mode"
// @Param file_order_method formData string false "file order method"
// @Param req body v1.CreateAuditTaskReqV1 true "create and audit task"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/projects/{project_name}/tasks/audits [post]
func CreateAndAuditTask(c echo.Context) error {
	// TODO 不同SQL模式审核增加备份配置
	req := new(CreateAuditTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	var sqls getSQLFromFileResp
	var err error
	var fileRecords []*model.AuditFile

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.Sql != "" {
		sqls = getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.Sql,
		}
	} else {
		sqls, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		fileRecords, err = saveFileFromContext(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	task, err := buildOnlineTaskForAudit(c, s, uint64(user.ID), req.InstanceName, req.InstanceSchema, projectUid, sqls)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// if task instance is not nil, gorm will update instance when save task.
	tmpInst := *task.Instance
	task.Instance = nil

	task.ExecMode = req.ExecMode
	task.FileOrderMethod = req.FileOrderMethod
	if req.EnableBackup {
		backupService := server.BackupService{}
		err = backupService.CheckBackupConflictWithExecMode(req.EnableBackup, req.ExecMode)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = backupService.CheckIsDbTypeSupportEnableBackup(task.DBType)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		task.EnableBackup = req.EnableBackup
	}
	task.InstanceEnableBackup = tmpInst.EnableBackup

	err = convertSQLSourceEncodingFromTask(task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskGroup := model.TaskGroup{Tasks: []*model.Task{task}}
	err = s.Save(&taskGroup)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(fileRecords) > 0 {
		fileHeader, _, err := getFileHeaderFromContext(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if strings.HasSuffix(fileHeader.Filename, ZIPFileExtension) && req.FileOrderMethod != "" && task.ExecMode == model.ExecModeSqlFile {
			sortAuditFiles(fileRecords, req.FileOrderMethod)
		}

		err = batchCreateFileRecords(s, fileRecords, task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.GenericError, fmt.Errorf("save sql file record failed: %v", err)))
		}
	}
	task.Instance = &tmpInst
	task, err = server.GetSqled().AddTaskWaitResult(projectUid, fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}

func convertSQLSourceEncodingFromTask(task *model.Task) error {
	for _, sql := range task.ExecuteSQLs {
		if sql.SourceFile == "" {
			continue
		}
		utf8NameByte, err := utils.ConvertToUtf8([]byte(sql.SourceFile))
		if err != nil {
			if e.Is(err, utils.ErrUnknownEncoding) {
				return e.New("the file name contains unrecognized characters. Please ensure the file name is encoded in UTF-8 or use an English file name")
			}
			return err
		}
		sql.SourceFile = string(utf8NameByte)
	}
	return nil
}

// @Summary 获取Sql扫描任务信息
// @Description get task
// @Tags task
// @Id getAuditTaskV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/tasks/audits/{task_id}/ [get]
func GetTask(c echo.Context) error {
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}
func GetTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	return getTaskById(ctx, taskId)
}

func getTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	s := model.GetStorage()
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.NewTaskNoExistOrNoAccessErr()
	}

	files, err := s.GetFileByTaskId(taskId)
	if err != nil {
		return nil, err
	}
	task.AuditFiles = files

	instance, exist, err := dms.GetInstancesById(ctx, fmt.Sprintf("%d", task.InstanceId))
	if err != nil {
		return nil, err
	}
	if exist {
		task.Instance = instance
	}

	return task, nil
}

type GetAuditTaskSQLsReqV1 struct {
	FilterExecStatus  string `json:"filter_exec_status" query:"filter_exec_status"`
	FilterAuditStatus string `json:"filter_audit_status" query:"filter_audit_status"`
	FilterAuditLevel  string `json:"filter_audit_level" query:"filter_audit_level"`
	NoDuplicate       bool   `json:"no_duplicate" query:"no_duplicate"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditTaskSQLsResV1 struct {
	controller.BaseRes
	Data      []*AuditTaskSQLResV1 `json:"data"`
	TotalNums uint64               `json:"total_nums"`
}

type AuditTaskSQLResV1 struct {
	Number       uint     `json:"number"`
	ExecSQL      string   `json:"exec_sql"`
	AuditResult  string   `json:"audit_result"`
	AuditLevel   string   `json:"audit_level"`
	AuditStatus  string   `json:"audit_status"`
	ExecResult   string   `json:"exec_result"`
	ExecStatus   string   `json:"exec_status"`
	RollbackSQLs []string `json:"rollback_sqls,omitempty"`
	Description  string   `json:"description"`
}

// @Summary 获取指定扫描任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed,manually_executed,execute_rollback)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param filter_audit_level query string false "filter: audit level of task sql" Enums(normal,notice,warn,error)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.GetAuditTaskSQLsResV1
// @router /v1/tasks/audits/{task_id}/sqls [get]
func GetTaskSQLs(c echo.Context) error {
	req := new(GetAuditTaskSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTaskDMS(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	data := map[string]interface{}{
		"task_id":             taskId,
		"filter_exec_status":  req.FilterExecStatus,
		"filter_audit_status": req.FilterAuditStatus,
		"filter_audit_level":  req.FilterAuditLevel,
		"no_duplicate":        req.NoDuplicate,
		"limit":               limit,
		"offset":              offset,
	}

	taskSQLs, count, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	rollbackSqlMap, err := server.BackupService{}.GetRollbackSqlsMap(task.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskSQLsRes := make([]*AuditTaskSQLResV1, 0, len(taskSQLs))
	for _, taskSQL := range taskSQLs {
		taskSQLRes := &AuditTaskSQLResV1{
			Number:       taskSQL.Number,
			Description:  taskSQL.Description,
			ExecSQL:      taskSQL.ExecSQL,
			AuditResult:  taskSQL.GetAuditResults(ctx),
			AuditLevel:   taskSQL.AuditLevel,
			AuditStatus:  taskSQL.AuditStatus,
			ExecResult:   taskSQL.ExecResult,
			ExecStatus:   taskSQL.ExecStatus,
			RollbackSQLs: rollbackSqlMap[taskSQL.Id],
		}
		taskSQLsRes = append(taskSQLsRes, taskSQLRes)
	}

	return c.JSON(http.StatusOK, &GetAuditTaskSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      taskSQLsRes,
		TotalNums: count,
	})
}

type DownloadAuditTaskSQLsFileReqV1 struct {
	NoDuplicate bool `json:"no_duplicate" query:"no_duplicate"`
}

// @Summary 下载指定扫描任务的SQLs信息报告
// @Description download report file of all SQLs information belong to the specified audit task
// @Tags task
// @Id downloadAuditTaskSQLReportV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Success 200 file 1 "sql report csv file"
// @router /v1/tasks/audits/{task_id}/sql_report [get]
func DownloadTaskSQLReportFile(c echo.Context) error {
	req := new(DownloadAuditTaskSQLsFileReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"task_id":      taskId,
		"no_duplicate": req.NoDuplicate,
	}

	taskSQLsDetail, _, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ctx := c.Request().Context()
	buff := &bytes.Buffer{}
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	cw := csv.NewWriter(buff)
	err = cw.Write([]string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportIndex),       // "序号",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportSQL),         // "SQL",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportAuditStatus), // "SQL审核状态",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportAuditResult), // "SQL审核结果",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportExecStatus),  // "SQL执行状态",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportExecResult),  // "SQL执行结果",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportRollbackSQL), // "SQL对应的回滚语句",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportDescription), // "SQL描述",
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.WriteDataToTheFileError, err))
	}
	rollbackSqlMap, err := server.BackupService{}.GetRollbackSqlsMap(task.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for _, td := range taskSQLsDetail {
		taskSql := &model.ExecuteSQL{
			AuditResults: td.AuditResults,
			AuditStatus:  td.AuditStatus,
		}
		taskSql.ExecStatus = td.ExecStatus
		err := cw.Write([]string{
			strconv.FormatUint(uint64(td.Number), 10),
			td.ExecSQL,
			taskSql.GetAuditStatusDesc(ctx),
			taskSql.GetAuditResultDesc(ctx),
			taskSql.GetExecStatusDesc(ctx),
			td.ExecResult,
			strings.Join(rollbackSqlMap[taskSql.ID], "\n"),
			td.Description,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.WriteDataToTheFileError, err))
		}
	}
	cw.Flush()
	fileName := fmt.Sprintf("SQL_audit_report_%v_%v.csv", task.InstanceName(), taskId)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))
	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

// @Summary 下载指定扫描任务的SQL文件
// @Description download SQL file for the audit task
// @Tags task
// @Id downloadAuditTaskSQLFileV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 file 1 "sql file"
// @router /v1/tasks/audits/{task_id}/sql_file [get]
func DownloadTaskSQLFile(c echo.Context) error {
	taskId := c.Param("task_id")
	s := model.GetStorage()
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	content, err := s.GetTaskExecuteSQLContent(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	fileName := fmt.Sprintf("exec_sql_%s_%s.sql", task.InstanceName(), taskId)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))

	return c.Blob(http.StatusOK, echo.MIMETextPlain, content)
}

type GetAuditTaskSQLContentResV1 struct {
	controller.BaseRes
	Data *AuditTaskSQLContentResV1 `json:"data"`
}

type AuditTaskSQLContentResV1 struct {
	Sql string `json:"sql" example:"alter table tb1 drop columns c1"`
}

// @Summary 获取指定扫描任务的SQL内容
// @Description get SQL content for the audit task
// @Tags task
// @Id getAuditTaskSQLContentV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetAuditTaskSQLContentResV1
// @router /v1/tasks/audits/{task_id}/sql_content [get]
func GetAuditTaskSQLContent(c echo.Context) error {
	taskId := c.Param("task_id")
	s := model.GetStorage()

	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	content, err := s.GetTaskExecuteSQLContent(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetAuditTaskSQLContentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &AuditTaskSQLContentResV1{
			Sql: string(content),
		},
	})
}

type UpdateAuditTaskSQLsReqV1 struct {
	Description string `json:"description"`
}

// @Summary 修改扫描任务中某条SQL的相关信息
// @Description modify the relevant information of a certain SQL in the audit task
// @Tags task
// @Id updateAuditTaskSQLsV1
// @Accept json
// @Param task_id path string true "task id"
// @Param number path string true "sql number"
// @Param audit_plan body v1.UpdateAuditTaskSQLsReqV1 true "modify the relevant information of a certain SQL in the audit task"
// @Security ApiKeyAuth
// @Success 200 {object} controller.BaseRes
// @router /v1/tasks/audits/{task_id}/sqls/{number} [patch]
func UpdateAuditTaskSQLs(c echo.Context) error {
	req := new(UpdateAuditTaskSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	taskId := c.Param("task_id")
	number := c.Param("number")

	s := model.GetStorage()
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanOpTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskSql, exist, err := s.GetTaskSQLByNumber(taskId, number)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("sql number not found")))
	}
	// the user may leave the description blank to clear the description, so no processing is performed
	taskSql.Description = req.Description
	err = s.Save(taskSql)
	return controller.JSONBaseErrorReq(c, err)
}

func CheckCurrentUserCanViewTask(c echo.Context, task *model.Task) (err error) {
	return checkCurrentUserCanViewTask(c, task, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeViewOthersWorkflow})
}

func CheckCurrentUserCanOpTask(c echo.Context, task *model.Task) (err error) {
	return checkCurrentUserCanOpTask(c, task, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeViewOthersWorkflow})
}

// TODO 使用DMS的权限校验
func CheckCurrentUserCanViewTaskDMS(c echo.Context, task *model.Task) error {
	_, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return err
	}
	return nil
}

type SQLExplain struct {
	SQL     string `json:"sql"`
	Message string `json:"message"`
	// explain result in table format
	ClassicResult ExplainClassicResult `json:"classic_result"`
}

type ExplainClassicResult struct {
	Rows []map[string] /* head name */ string `json:"rows"`
	Head []TableMetaItemHeadResV1             `json:"head"`
}

type GetTaskAnalysisDataResItemV1 struct {
	SQLExplain SQLExplain  `json:"sql_explain"`
	TableMetas []TableMeta `json:"table_metas"`
}

type GetTaskAnalysisDataResV1 struct {
	controller.BaseRes
	Data GetTaskAnalysisDataResItemV1 `json:"data"`
}

// GetTaskAnalysisData get SQL explain and related table metadata for analysis
// @Summary 获取task相关的SQL执行计划和表元数据
// @Description get SQL explain and related table metadata for analysis
// @Id getTaskAnalysisData
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskAnalysisDataResV1
// @router /v1/tasks/audits/{task_id}/sqls/{number}/analysis [get]
func GetTaskAnalysisData(c echo.Context) error {
	return getTaskAnalysisData(c)
}

type CreateAuditTasksGroupReqV1 struct {
	Instances       []*InstanceForCreatingTask `json:"instances" valid:"dive,required"`
	ExecMode        string                     `json:"exec_mode" enums:"sql_file,sqls"`
	FileOrderMethod string                     `json:"file_order_method"`
}

type InstanceForCreatingTask struct {
	InstanceName   string `json:"instance_name" form:"instance_name" valid:"required"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema"`
}

type CreateAuditTasksGroupResV1 struct {
	controller.BaseRes
	Data AuditTasksGroupResV1 `json:"data"`
}

type AuditTasksGroupResV1 struct {
	TaskGroupId uint `json:"task_group_id" form:"task_group_id" valid:"required"`
}

const MaximumDataSourceNum = 10

// CreateAuditTasksGroupV1
// @Summary 创建审核任务组
// @Description create tasks group.
// @Accept json
// @Produce json
// @Tags task
// @Id createAuditTasksV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param req body v1.CreateAuditTasksGroupReqV1 true "parameters for creating audit tasks group"
// @Success 200 {object} v1.CreateAuditTasksGroupResV1
// @router /v1/projects/{project_name}/task_groups [post]
func CreateAuditTasksGroupV1(c echo.Context) error {
	req := new(CreateAuditTasksGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 数据源个数最大为10
	if len(req.Instances) > MaximumDataSourceNum {
		return controller.JSONBaseErrorReq(c, ErrTooManyDataSource)
	}

	instNames := make([]string, len(req.Instances))
	for i, instance := range req.Instances {
		instNames[i] = instance.InstanceName
	}

	distinctInstNames := utils.RemoveDuplicate(instNames)

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	instances, err := dms.GetInstancesInProjectByNames(c.Request().Context(), projectUid, distinctInstNames)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	nameInstanceMap := make(map[string]*model.Instance, len(req.Instances))
	for _, inst := range instances {
		// https://github.com/actiontech/sqle/issues/1673
		inst, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, inst.Name)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
		}

		can, err := CheckCurrentUserCanOpInstances(c.Request().Context(), projectUid, user.GetIDStr(), instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !can {
			return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
		}

		if err := common.CheckInstanceIsConnectable(inst); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		nameInstanceMap[inst.Name] = inst
	}

	tasks := make([]*model.Task, len(req.Instances))
	for i, reqInstance := range req.Instances {
		tasks[i] = &model.Task{
			Schema:       reqInstance.InstanceSchema,
			InstanceId:   nameInstanceMap[reqInstance.InstanceName].ID,
			CreateUserId: uint64(user.ID),
			DBType:       nameInstanceMap[reqInstance.InstanceName].DbType,
		}
		tasks[i].CreatedAt = time.Now()
		tasks[i].ExecMode = req.ExecMode
		tasks[i].FileOrderMethod = req.FileOrderMethod
	}

	taskGroup := model.TaskGroup{Tasks: tasks}
	if err := s.Save(&taskGroup); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, CreateAuditTasksGroupResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditTasksGroupResV1{
			TaskGroupId: taskGroup.ID,
		},
	})
}

type AuditTaskGroupReqV1 struct {
	TaskGroupId  uint   `json:"task_group_id" form:"task_group_id" valid:"required"`
	Sql          string `json:"sql" form:"sql" example:"alter table tb1 drop columns c1"`
	EnableBackup bool   `json:"enable_backup" form:"enable_backup"`
}

type AuditTaskGroupRes struct {
	TaskGroupId uint              `json:"task_group_id"`
	Tasks       []*AuditTaskResV1 `json:"tasks"`
}

type AuditTaskGroupResV1 struct {
	controller.BaseRes
	Data AuditTaskGroupRes `json:"data"`
}

// AuditTaskGroupV1
// @Summary 审核任务组
// @Description audit task group.
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is zip file, sql will be parsed from it.
// @Accept mpfd
// @Produce json
// @Tags task
// @Id auditTaskGroupIdV1
// @Security ApiKeyAuth
// @Param task_group_id formData uint true "group id of tasks"
// @Param sql formData string false "sqls for audit"
// @Param enable_backup formData bool false "enable backup"
// @Param file_order_method formData string false "file order method"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.AuditTaskGroupResV1
// @router /v1/task_groups/audit [post]
func AuditTaskGroupV1(c echo.Context) error {
	// TODO 单数据源审核，以及多数据源相同SQL模式审核，增加备份配置
	req := new(AuditTaskGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	var err error
	var sqls getSQLFromFileResp
	var fileRecords []*model.AuditFile
	s := model.GetStorage()
	taskGroup, err := s.GetTaskGroupByGroupId(req.TaskGroupId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	tasks := taskGroup.Tasks
	instanceIds := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIds = append(instanceIds, task.InstanceId)
	}

	instances, err := dms.GetInstancesByIds(c.Request().Context(), instanceIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 因为这个接口数据源属于同一个项目,取第一个DB所属项目
	projectId := instances[0].ProjectId
	can, err := CheckCurrentUserCanOpInstances(c.Request().Context(), projectId, controller.GetUserID(c), instances)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	// 因为这个接口数据源必然相同，所以只取第一个实例的DbType即可
	dbType := instances[0].DbType

	if req.Sql != "" {
		sqls = getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.Sql,
		}
	} else {
		sqls, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		fileRecords, err = saveFileFromContext(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	{
		l := log.NewEntry()

		plugin, err := common.NewDriverManagerWithoutCfg(l, dbType)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		defer plugin.Close(context.TODO())
		instanceMap := make(map[uint64]*model.Instance)
		for _, instance := range instances {
			instanceMap[instance.ID] = instance
		}

		for _, task := range tasks {
			task.SQLSource = sqls.SourceType
			if req.EnableBackup {
				backupService := server.BackupService{}
				err = backupService.CheckBackupConflictWithExecMode(req.EnableBackup, task.ExecMode)
				if err != nil {
					return controller.JSONBaseErrorReq(c, err)
				}
				err = backupService.CheckIsDbTypeSupportEnableBackup(task.DBType)
				if err != nil {
					return controller.JSONBaseErrorReq(c, err)
				}
				task.EnableBackup = req.EnableBackup
			}
			if instance, exist := instanceMap[task.InstanceId]; exist {
				task.InstanceEnableBackup = instance.EnableBackup
			} else {
				return controller.JSONBaseErrorReq(c, fmt.Errorf("can not find instance in task"))
			}
			err := addSQLsFromFileToTasks(sqls, task, plugin)
			if err != nil {
				return controller.JSONBaseErrorReq(c, errors.New(errors.GenericError, fmt.Errorf("add sqls from file to task failed: %v", err)))
			}
			if len(fileRecords) > 0 {
				fileHeader, _, err := getFileHeaderFromContext(c)
				if err != nil {
					return controller.JSONBaseErrorReq(c, err)
				}
				if strings.HasSuffix(fileHeader.Filename, ZIPFileExtension) && task.FileOrderMethod != "" && task.ExecMode == model.ExecModeSqlFile {
					sortAuditFiles(fileRecords, task.FileOrderMethod)
				}

				err = batchCreateFileRecords(s, fileRecords, task.ID)
				if err != nil {
					return controller.JSONBaseErrorReq(c, errors.New(errors.GenericError, fmt.Errorf("save sql file record failed: %v", err)))
				}
			}
		}
	}

	for _, task := range taskGroup.Tasks {
		err = convertSQLSourceEncodingFromTask(task)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if err := s.Save(taskGroup.Tasks); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	for i, task := range tasks {
		if task.Status != model.TaskStatusInit {
			continue
		}

		tasks[i], err = server.GetSqled().AddTaskWaitResult(projectId, fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	tasksRes := make([]*AuditTaskResV1, len(tasks))
	for i, task := range tasks {
		tasksRes[i] = &AuditTaskResV1{
			Id:                         task.ID,
			InstanceName:               task.InstanceName(),
			InstanceDbType:             task.DBType,
			InstanceSchema:             task.Schema,
			AuditLevel:                 task.AuditLevel,
			Score:                      task.Score,
			PassRate:                   task.PassRate,
			Status:                     task.Status,
			EnableBackup:               task.EnableBackup,
			BackupConflictWithInstance: server.BackupService{}.IsBackupConflictWithInstance(task.EnableBackup, task.InstanceEnableBackup),
			SQLSource:                  task.SQLSource,
			ExecStartTime:              task.ExecStartAt,
			ExecEndTime:                task.ExecEndAt,
		}
	}

	return c.JSON(http.StatusOK, AuditTaskGroupResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditTaskGroupRes{
			TaskGroupId: taskGroup.ID,
			Tasks:       tasksRes,
		},
	})
}

func batchCreateFileRecords(s *model.Storage, fileRecords []*model.AuditFile, taskId uint) error {
	// Initialize parentID to 0
	var parentID uint
	// Create a slice to store all file records except the first one
	records := make([]*model.AuditFile, 0, len(fileRecords)-1)

	for i, fileRecord := range fileRecords {
		// Set TaskId and ParentID for each file record
		fileRecord.TaskId = taskId
		fileRecord.ParentID = parentID

		if i == 0 {
			fileRecord.ID = 0
			// save first record as the parent file record
			if err := s.Create(fileRecord); err != nil {
				return err
			}
			// Update parentID to the ID of the first record
			parentID = fileRecord.ID
		} else {
			// Add the record to records slice
			records = append(records, fileRecord)
		}
	}
	// Batch save all file records except the first one
	if len(records) > 0 {
		if err := s.BatchCreateFileRecords(records); err != nil {
			return err
		}
	}
	return nil
}

// @Summary 获取指定审核任务的原始文件
// @Description get SQL origin file of the audit task
// @Tags task
// @Id DownloadAuditFile
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} controller.BaseRes
// @router /v1/tasks/audits/{task_id}/origin_file [get]
func DownloadAuditFile(c echo.Context) error {
	taskId := c.Param("task_id")
	s := model.GetStorage()
	files, err := s.GetFileByTaskId(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	log.Logger().Debugf("task id %v", taskId)
	/*
		TODO 鉴权
		err = CheckCurrentUserCanViewTask(c, task)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	*/
	if len(files) == 0 {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("can not find any file in this task"))
	}

	file := files[0]
	if file.FileHost != config.GetOptions().SqleOptions.ReportHost {
		log.NewEntry().Debugf("try to reverse to sqle due to file.FileHost %v this host %v", file.FileHost, config.GetOptions().SqleOptions.ReportHost)
		err = ReverseToSqle(c, c.Request().URL.Path, file.FileHost)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		filePath := model.DefaultFilePath(file.UniqueName)
		err = c.Attachment(filePath, file.FileName)
		if err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusOK)

}

// TODO 和DMS一起抽离出一个工具函数
func ReverseToSqle(c echo.Context, rewriteUrlPath, targetHost string) (err error) {
	// c.Request().URL.Path = rewriteUrlPath
	// reference from echo framework proxy middleware
	target, err := url.Parse(fmt.Sprintf("http://%s", targetHost))
	log.NewEntry().Debugf("reverse to sqle: %v%v", target.Host, c.Request().URL.Path)
	if err != nil {
		return err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		// If the client canceled the request (usually by closing the connection), we can report a
		// client error (4xx) instead of a server error (5xx) to correctly identify the situation.
		// The Go standard library (at of late 2020) wraps the exported, standard
		// context.Canceled error with unexported garbage value requiring a substring check, see
		// https://github.com/golang/go/blob/6965b01ea248cabb70c3749fd218b36089a21efb/src/net/net.go#L416-L430
		if err == context.Canceled || strings.Contains(err.Error(), "operation was canceled") {
			httpError := echo.NewHTTPError(middleware.StatusCodeContextCanceled, fmt.Sprintf("client closed connection: %v", err))
			httpError.Internal = err
			c.Set("_error", httpError)
		} else {
			httpError := echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("remote %s unreachable, could not forward: %v", target.String(), err))
			httpError.Internal = err
			c.Set("_error", httpError)
		}
	}

	reverseProxy.ServeHTTP(c.Response(), c.Request())

	if e, ok := c.Get("_error").(error); ok {
		err = e
	}

	return
}

type SqlFileOrderMethod struct {
	OrderMethod string `json:"order_method"`
	Desc        string `json:"desc"`
}

type SqlFileOrderMethodRes struct {
	Methods []SqlFileOrderMethod `json:"methods"`
}

type GetSqlFileOrderMethodResV1 struct {
	controller.BaseRes
	Data SqlFileOrderMethodRes `json:"data"`
}

// GetSqlFileOrderMethodV1
// @Summary 获取文件上线排序方式
// @Description get file order method
// @Accept json
// @Produce json
// @Tags task
// @Id getSqlFileOrderMethodV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSqlFileOrderMethodResV1
// @router /v1/tasks/file_order_methods [get]
func GetSqlFileOrderMethodV1(c echo.Context) error {
	return getSqlFileOrderMethod(c)
}
