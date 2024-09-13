//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/labstack/echo/v4"
)

func getTaskAnalysisData(c echo.Context) error {
	taskId := c.Param("task_id")
	number := c.Param("number")

	s := model.GetStorage()
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewTaskNoExistOrNoAccessErr())
	}
	instance, exist, err := dms.GetInstancesById(c.Request().Context(), fmt.Sprintf("%d", task.InstanceId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewTaskNoExistOrNoAccessErr())
	}

	task.Instance = instance

	if err := CheckCurrentUserCanViewTask(c, task); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskSql, exist, err := s.GetTaskSQLByNumber(taskId, number)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("sql number not found")))
	}

	explainResult, explainMessage, metaDataResult, err := getSQLAnalysisResultFromDriver(log.NewEntry(), task.Schema, taskSql.Content, task.Instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTaskAnalysisDataResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertExplainAndMetaDataToRes(c.Request().Context(), explainResult, explainMessage, metaDataResult, taskSql.Content),
	})
}

func convertExplainAndMetaDataToRes(ctx context.Context, explainResultInput *driverV2.ExplainResult, explainMessage string, metaDataResultInput *driver.GetTableMetaBySQLResult,
	rawSql string) GetTaskAnalysisDataResItemV1 {

	explainResult := explainResultInput
	if explainResult == nil {
		explainResult = &driverV2.ExplainResult{}
	}
	metaDataResult := metaDataResultInput
	if metaDataResult == nil {
		metaDataResult = &driver.GetTableMetaBySQLResult{}
	}

	analysisDataResItemV1 := GetTaskAnalysisDataResItemV1{
		SQLExplain: SQLExplain{
			ClassicResult: ExplainClassicResult{
				Rows: make([]map[string]string, len(explainResult.ClassicResult.Rows)),
				Head: make([]TableMetaItemHeadResV1, len(explainResult.ClassicResult.Columns)),
			},
		},
		TableMetas: make([]TableMeta, len(metaDataResult.TableMetas)),
	}

	explainResItemV1 := analysisDataResItemV1.SQLExplain.ClassicResult
	for i, column := range explainResult.ClassicResult.Columns {
		explainResItemV1.Head[i].FieldName = column.Name
		explainResItemV1.Head[i].Desc = column.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
	}

	for i, rows := range explainResult.ClassicResult.Rows {
		explainResItemV1.Rows[i] = make(map[string]string)
		for k, row := range rows {
			columnName := explainResult.ClassicResult.Columns[k].Name
			explainResItemV1.Rows[i][columnName] = row
		}
	}

	analysisDataResItemV1.SQLExplain.SQL = rawSql
	analysisDataResItemV1.SQLExplain.Message = explainMessage

	for i, tableMeta := range metaDataResult.TableMetas {
		tableMetaColumnsInfo := tableMeta.ColumnsInfo
		tableMetaIndexInfo := tableMeta.IndexesInfo

		analysisDataResItemV1.TableMetas[i].Columns = TableColumns{
			Rows: make([]map[string]string, len(tableMetaColumnsInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(tableMetaColumnsInfo.Columns)),
		}
		analysisDataResItemV1.TableMetas[i].Indexes = TableIndexes{
			Rows: make([]map[string]string, len(tableMetaIndexInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(tableMetaIndexInfo.Columns)),
		}

		tableMetaColumnRes := analysisDataResItemV1.TableMetas[i].Columns
		for i2, column := range tableMetaColumnsInfo.Columns {
			tableMetaColumnRes.Head[i2].FieldName = column.Name
			tableMetaColumnRes.Head[i2].Desc = column.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
		}

		for i2, rows := range tableMetaColumnsInfo.Rows {
			tableMetaColumnRes.Rows[i2] = make(map[string]string)
			for k, row := range rows {
				columnName := tableMetaColumnsInfo.Columns[k].Name
				tableMetaColumnRes.Rows[i2][columnName] = row
			}
		}

		tableMetaIndexRes := analysisDataResItemV1.TableMetas[i].Indexes
		for i2, column := range tableMetaIndexInfo.Columns {
			tableMetaIndexRes.Head[i2].FieldName = column.Name
			tableMetaIndexRes.Head[i2].Desc = column.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
		}

		for i2, rows := range tableMetaIndexInfo.Rows {
			tableMetaIndexRes.Rows[i2] = make(map[string]string)
			for k, row := range rows {
				columnName := tableMetaIndexInfo.Columns[k].Name
				tableMetaIndexRes.Rows[i2][columnName] = row
			}
		}

		analysisDataResItemV1.TableMetas[i].Name = tableMeta.Name
		analysisDataResItemV1.TableMetas[i].Schema = tableMeta.Schema
		analysisDataResItemV1.TableMetas[i].CreateTableSQL = tableMeta.CreateTableSQL
		analysisDataResItemV1.TableMetas[i].Message = tableMeta.Message
	}

	return analysisDataResItemV1
}

const (
	FileOrderMethodPrefixNumAsc = "file_order_method_prefix_num_asc"
	FileOrderMethodSuffixNumAsc = "file_order_method_suffix_num_asc"
)

type FileOrderMethod struct {
	Method string
	Desc   *i18n.Message
}

var FileOrderMethods = []FileOrderMethod{
	{
		Method: FileOrderMethodPrefixNumAsc,
		Desc:   locale.FileOrderMethodPrefixNumAsc,
	},
	{
		Method: FileOrderMethodSuffixNumAsc,
		Desc:   locale.FileOrderMethodSuffixNumAsc,
	},
}

func getSqlFileOrderMethod(c echo.Context) error {
	methods := make([]SqlFileOrderMethod, 0, len(FileOrderMethods))
	for _, method := range FileOrderMethods {
		methods = append(methods, SqlFileOrderMethod{
			OrderMethod: method.Method,
			Desc:        locale.ShouldLocalizeMsg(c.Request().Context(), method.Desc),
		})
	}
	return c.JSON(http.StatusOK, GetSqlFileOrderMethodResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SqlFileOrderMethodRes{
			Methods: methods,
		},
	})
}

type auditFileWithNum struct {
	auditFile *model.AuditFile
	num       int
}

type auditFileWithNums []auditFileWithNum

func (s auditFileWithNums) Len() int           { return len(s) }
func (s auditFileWithNums) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s auditFileWithNums) Less(i, j int) bool { return s[i].num < s[j].num }

func sortAuditFiles(auditFiles []*model.AuditFile, orderMethod string) {
	var re *regexp.Regexp
	sortedAuditFiles := []*model.AuditFile{}

	if len(auditFiles) <= 1 {
		return
	}

	// 索引为0的auditFiles为zip包自身的信息，不参与排序
	sortedAuditFiles = append(sortedAuditFiles, auditFiles[0])
	auditFiles = auditFiles[1:]

	switch orderMethod {
	case FileOrderMethodPrefixNumAsc:
		re = regexp.MustCompile(`^\d+`)
	case FileOrderMethodSuffixNumAsc:
		re = regexp.MustCompile(`\d+$`)
	}
	if re == nil {
		return
	}

	fileWithNums, invalidOrderFiles := getFileWithNumFromPathsByRegexp(auditFiles, re)
	fileWithSortNums := auditFileWithNums(fileWithNums)
	sort.Sort(fileWithSortNums)

	for _, fileWithSortNum := range fileWithSortNums {
		sortedAuditFiles = append(sortedAuditFiles, fileWithSortNum.auditFile)
	}
	sortedAuditFiles = append(sortedAuditFiles, invalidOrderFiles...)
	for i, auditFile := range sortedAuditFiles {
		auditFile.ExecOrder = uint(i)
	}
}

func getFileWithNumFromPathsByRegexp(auditFiles []*model.AuditFile, re *regexp.Regexp) ([]auditFileWithNum, []*model.AuditFile) {
	invalidOrderFiles := []*model.AuditFile{} // 不符合排序规则的文件路径
	fileWithNums := []auditFileWithNum{}

	for _, file := range auditFiles {
		filename := getFileNameWithoutExtension(file.FileName)
		match := re.FindString(filename)
		if match == "" {
			invalidOrderFiles = append(invalidOrderFiles, file)
			log.NewEntry().Errorf("getSortNumsFromFilePaths regexp match failed, filename:%s, regexp:%s", file.FileName, re.String())
			continue
		}
		num, err := strconv.Atoi(match)
		if err != nil {
			invalidOrderFiles = append(invalidOrderFiles, file)
			log.NewEntry().Errorf("getSortNumsFromFilePaths convert string to number failed, string:%s, filename:%s,  err:%v", match, file.FileName, err)
			continue
		}
		fileWithNums = append(fileWithNums, auditFileWithNum{
			auditFile: file,
			num:       num,
		})
	}
	return fileWithNums, invalidOrderFiles
}

func getFileNameWithoutExtension(filePath string) string {
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}
