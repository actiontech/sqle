package auditplan

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/sirupsen/logrus"
)

var errNoSQLInAuditPlan = errors.New(errors.DataConflict, fmt.Errorf("there is no SQLs in audit plan"))
var errNoSQLNeedToBeAudited = errors.New(errors.DataConflict, fmt.Errorf("there is no SQLs need to be audited in audit plan"))

type Task interface {
	Start() error
	Stop() error
	Audit() (*AuditResultResp, error)
	FullSyncSQLs([]*SQL) error
	PartialSyncSQLs([]*SQL) error
	GetSQLs(map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error)
}

func NewTask(entry *logrus.Entry, ap *AuditPlan) Task {
	entry = entry.WithField("id", ap.ID)

	meta, err := GetMeta(ap.Type)
	if err != nil || meta.CreateTask == nil {
		return NewDefaultTask(entry, ap)
	}

	return meta.CreateTask(entry, ap)
}

type baseTask struct {
	ap *AuditPlan
	// persist is a database handle which store AuditPlan.
	persist *model.Storage
	logger  *logrus.Entry
}

func newBaseTask(entry *logrus.Entry, ap *AuditPlan) *baseTask {
	log.NewEntry()
	return &baseTask{
		ap:      ap,
		persist: model.GetStorage(),
		logger:  entry,
	}
}

func (at *baseTask) Start() error {
	return nil
}

func (at *baseTask) Stop() error {
	return nil
}

type AuditPlanSQL struct {
	ID             uint
	AuditPlanID    uint       `json:"audit_plan_id"`
	Fingerprint    string     `json:"fingerprint"`
	FingerprintMD5 string     `json:"fingerprint_md5"`
	SQLContent     string     `json:"sql"`
	Info           model.JSON `gorm:"type:json"`
	Schema         string     `json:"schema"`
}

func ConvertModel2AuditPlanSql(m model.AuditPlanSQLV2) AuditPlanSQL {
	return AuditPlanSQL{
		ID:             m.ID,
		AuditPlanID:    m.AuditPlanID,
		Fingerprint:    m.Fingerprint,
		FingerprintMD5: m.FingerprintMD5,
		SQLContent:     m.SQLContent,
		Info:           m.Info,
		Schema:         m.Schema,
	}
}

func (at *baseTask) audit(task *model.Task) (*AuditResultResp, error) {
	auditPlanSQLs, err := at.persist.GetAuditPlanSQLsV2Unaudit(at.ap.ID)
	if err != nil {
		return nil, err
	}

	if len(auditPlanSQLs) == 0 {
		return nil, errNoSQLInAuditPlan
	}

	filteredSqls, err := filterSQLsByPeriodV2(at.ap.Params, auditPlanSQLs)
	if err != nil {
		return nil, err
	}

	if len(filteredSqls) == 0 {
		return nil, errNoSQLNeedToBeAudited
	}

	for i, sql := range filteredSqls {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SqlText,
			},
		})
	}
	projectId := model.ProjectUID(at.ap.ProjectId)
	err = server.Audit(at.logger, task, &projectId, at.ap.RuleTemplateName)
	if err != nil {
		return nil, err
	}

	return &AuditResultResp{
		AuditPlanID:  uint64(at.ap.ID),
		Task:         task,
		FilteredSqls: filteredSqls,
	}, nil
}

func filterSQLsByPeriodV2(params params.Params, sqls []*model.OriginManageSQL) (filteredSqls []*model.OriginManageSQL, err error) {
	period := params.GetParam(paramKeyAuditSQLsScrappedInLastPeriodMinute).Int()
	if period <= 0 {
		return sqls, nil
	}

	t := time.Now()
	minus := -1
	startTime := t.Add(time.Minute * time.Duration(minus*period))
	for _, sql := range sqls {
		var info = struct {
			LastReceiveTimestamp time.Time `json:"last_receive_timestamp"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, fmt.Errorf("parse last_receive_timestamp failed: %v", err)
		}

		if info.LastReceiveTimestamp.Before(startTime) {
			continue
		}
		newSql := *sql
		filteredSqls = append(filteredSqls, &newSql)
	}
	return filteredSqls, nil
}

// func filterSQLsByPeriod(params params.Params, sqls []*model.AuditPlanSQLV2) (filteredSqls []*model.AuditPlanSQLV2, err error) {
// 	period := params.GetParam(paramKeyAuditSQLsScrappedInLastPeriodMinute).Int()
// 	if period <= 0 {
// 		return sqls, nil
// 	}

// 	t := time.Now()
// 	minus := -1
// 	startTime := t.Add(time.Minute * time.Duration(minus*period))
// 	for _, sql := range sqls {
// 		var info = struct {
// 			LastReceiveTimestamp time.Time `json:"last_receive_timestamp"`
// 		}{}
// 		err := json.Unmarshal(sql.Info, &info)
// 		if err != nil {
// 			return nil, fmt.Errorf("parse last_receive_timestamp failed: %v", err)
// 		}

// 		if info.LastReceiveTimestamp.Before(startTime) {
// 			continue
// 		}
// 		newSql := *sql
// 		filteredSqls = append(filteredSqls, &newSql)
// 	}
// 	return filteredSqls, nil
// }

type sqlCollector struct {
	*baseTask
	sync.WaitGroup
	isStarted bool
	cancel    chan struct{}
	do        func()

	loopInterval func() time.Duration
}

// func newSQLCollector(entry *logrus.Entry, ap *AuditPlan) *sqlCollector {
// 	return &sqlCollector{
// 		newBaseTask(entry, ap),
// 		sync.WaitGroup{},
// 		false,
// 		make(chan struct{}),
// 		func() { // default
// 			entry.Warn("sql collector do nothing")
// 		},
// 		func() time.Duration {
// 			interval := ap.Params.GetParam(paramKeyCollectIntervalMinute).Int()
// 			if interval == 0 {
// 				interval = 60
// 			}
// 			return time.Minute * time.Duration(interval)
// 		},
// 	}
// }

func (at *sqlCollector) Start() error {
	if at.isStarted {
		return nil
	}
	interval := at.loopInterval()

	at.WaitGroup.Add(1)
	go func() {
		at.isStarted = true
		at.logger.Infof("start task")
		at.loop(at.cancel, interval)
		at.WaitGroup.Done()
	}()
	return nil
}

func (at *sqlCollector) Stop() error {
	if !at.isStarted {
		return nil
	}
	at.cancel <- struct{}{}
	at.isStarted = false
	at.WaitGroup.Wait()
	at.logger.Infof("stop task")
	return nil
}

func (at *sqlCollector) FullSyncSQLs(sqls []*SQL) error {
	at.logger.Warnf("someone try to sync sql to audit plan(%v), but sql should collected by task itself", at.ap.Name)
	return nil
}

func (at *sqlCollector) PartialSyncSQLs(sqls []*SQL) error {
	at.logger.Warnf("someone try to sync sql to audit plan(%v), but sql should collected by task itself", at.ap.Name)
	return nil
}

func (at *sqlCollector) loop(cancel chan struct{}, interval time.Duration) {
	at.do()
	if interval == 0 {
		at.logger.Warnf("task(%v) loop interval can not be zero", at.ap.Name)
		return
	}

	tk := time.NewTicker(interval)
	for {
		select {
		case <-cancel:
			tk.Stop()
			return
		case <-tk.C:
			at.logger.Infof("tick %s", at.ap.Name)
			at.do()
		}
	}
}

type DefaultTask struct {
	*baseTask
}

func NewDefaultTask(entry *logrus.Entry, ap *AuditPlan) Task {
	return &DefaultTask{newBaseTask(entry, ap)}
}

func (at *DefaultTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
	}
	return at.baseTask.audit(task)
}

func convertSQLsToModelSQLs(sqls []*SQL) []*model.AuditPlanSQLV2 {
	as := make([]*model.AuditPlanSQLV2, len(sqls))
	for i, sql := range sqls {
		data, _ := json.Marshal(sql.Info)
		as[i] = &model.AuditPlanSQLV2{
			Fingerprint: sql.Fingerprint,
			SQLContent:  sql.SQLContent,
			Schema:      sql.Schema,
			Info:        data,
		}
	}
	return as
}

func (at *baseTask) FullSyncSQLs(sqls []*SQL) error {
	return at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *baseTask) PartialSyncSQLs(sqls []*SQL) error {
	return at.persist.UpdateDefaultAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *baseTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	return baseTaskGetSQLs(args, at.persist)
}

func baseTaskGetSQLs(args map[string]interface{}, persist *model.Storage) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "fingerprint",
			Desc: "SQL指纹",
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: "最后一次匹配到该指纹的语句",
			Type: "sql",
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: "counter",
			Desc: "匹配到该指纹的语句数量",
		},
		{
			Name: "last_receive_timestamp",
			Desc: "最后一次匹配到该指纹的时间",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		var info = struct {
			Counter              uint64 `json:"counter"`
			LastReceiveTimestamp string `json:"last_receive_timestamp"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, nil, 0, err
		}
		rows = append(rows, map[string]string{
			"sql":                    sql.SQLContent,
			"fingerprint":            sql.Fingerprint,
			"counter":                strconv.FormatUint(info.Counter, 10),
			"last_receive_timestamp": info.LastReceiveTimestamp,
			model.AuditResultName:    sql.AuditResult.String,
		})
	}
	return head, rows, count, nil
}

func getTaskWithInstanceByAuditPlan(ap *AuditPlan, persist *model.Storage) (*model.Task, error) {
	var task *model.Task
	if ap.InstanceName == "" {
		task = &model.Task{
			DBType: ap.DBType,
		}
	} else {
		instance, _, err := dms.GetInstanceInProjectByName(context.TODO(), ap.ProjectId, ap.InstanceName)
		if err != nil {
			return nil, err
		}
		task = &model.Task{
			Instance: instance,
			DBType:   ap.DBType,
		}
	}
	return task, nil
}

// type sqlInfo struct {
// 	counter          int
// 	fingerprint      string
// 	sql              string
// 	schema           string
// 	queryTimeSeconds int
// 	//nolint:unused
// 	startTime string
// 	//nolint:unused
// 	rowExaminedAvg float64
// }

// func mergeSQLsByFingerprint(sqls []SqlFromAliCloud) []sqlInfo {
// 	sqlInfos := []sqlInfo{}

// 	counter := map[string]int /*slice subscript*/ {}
// 	for _, sql := range sqls {
// 		fp := query.Fingerprint(sql.sql)
// 		if index, exist := counter[fp]; exist {
// 			sqlInfos[index].counter += 1
// 			sqlInfos[index].fingerprint = fp
// 			sqlInfos[index].sql = sql.sql
// 			sqlInfos[index].schema = sql.schema
// 		} else {
// 			sqlInfos = append(sqlInfos, sqlInfo{
// 				counter:     1,
// 				fingerprint: fp,
// 				sql:         sql.sql,
// 				schema:      sql.schema,
// 			})
// 			counter[fp] = len(sqlInfos) - 1
// 		}

// 	}
// 	return sqlInfos
// }
