//go:build enterprise
// +build enterprise

package tbase_audit_log

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/actiontech/sqle/sqle/utils"
	pgParser "github.com/pganalyze/pg_query_go/v4/parser"

	"github.com/sirupsen/logrus"
)

type AuditLog struct {
	l *logrus.Entry
	c *scanner.Client

	logFolder string

	sqlCh chan scanners.SQL

	auditPlanID string
	parser      *TbaseParser
	Watcher     *Watcher
	fileFormat  string
}

type Params struct {
	LogFolder      string
	AuditPlanID    string
	FileNameFormat string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*AuditLog, error) {
	return &AuditLog{
		l:           l,
		c:           c,
		sqlCh:       make(chan scanners.SQL, 10),
		Watcher:     NewWatcher(),
		parser:      NewTbaseParser(),
		logFolder:   params.LogFolder,
		auditPlanID: params.AuditPlanID,
		fileFormat:  params.FileNameFormat,
	}, nil
}

func (a *AuditLog) GetSQLFromCsvFile(ctx context.Context, filePath string, waitFileUpdate bool) {
	file, err := os.Open(filePath)
	if err != nil {
		a.l.Errorf("tbase log open file failed, error: %v", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			if !waitFileUpdate {
				break
			}
			// 如果 waitFileUpdate 为 true，表示读取的是最新的日志文件，
			// 在遇到 EOF 时，不立即退出，而是等待新的日志内容写入。
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(5 * time.Second)
				continue
			}
		}
		csvError, ok := err.(*csv.ParseError)
		if ok {
			if csvError.Err == csv.ErrFieldCount {
				continue
			}
		}
		if err != nil {
			a.l.Errorf("tbase log read csv record failed, error: %v", err)
			return
		}
		tlog, err := a.parser.parseRecord(record)
		if err != nil {
			a.l.Errorf("tbase log parse record failed, error: %v", err)
			continue
		}
		if tlog != nil {
			a.sqlCh <- scanners.SQL{
				RawText:   tlog.SQLText,
				Counter:   1,
				Schema:    tlog.Schema,
				QueryTime: tlog.Duration,
				QueryAt:   tlog.ConnectionTime,
				DBUser:    tlog.User,
				Endpoint:  tlog.ClientHostWithPort,
			}
		}
	}
}

// 客户日志的保留时间为7天,日志默认每两小时切换一个
// 所以需要解析现有的日志文件，还需要监听是否有新的日志文件被创建，如果有，需要读取新的日志文件
// 具体参见：https://github.com/actiontech/sqle-ee/issues/1621
func (t *AuditLog) Run(ctx context.Context) error {
	sortedFiles, err := getSortedFilesByFolderPath(path.Join(t.logFolder, t.fileFormat))
	if err != nil {
		return err
	}
	newestLogFile := sortedFiles[0]
	go func() {
		for _, path := range sortedFiles[1:] {
			t.GetSQLFromCsvFile(context.Background(), path, false)
		}
	}()

	go t.Watcher.WatchFileCreated(filepath.Dir(newestLogFile), t.fileFormat)

	ctx, cancel := context.WithCancel(context.Background())
	go t.GetSQLFromCsvFile(ctx, newestLogFile, true)

	for {
		select {
		case <-ctx.Done():
			close(t.sqlCh)
			cancel()
			return nil

		case f, ok := <-t.Watcher.NewFiles:
			if !ok {
				cancel()
				return nil
			}

			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			go t.GetSQLFromCsvFile(ctx, f, true)
		}
	}
}

func (t *AuditLog) SQLs() <-chan scanners.SQL {
	return t.sqlCh
}

type sqlIdentity struct {
	Fingerprint string
	Schema      string
}

func (s sqlIdentity) jsonString() string {
	keyByte, _ := json.Marshal(s)
	return string(keyByte)
}

func (t *AuditLog) Upload(ctx context.Context, sqls []scanners.SQL) error {
	var sqlListReq []*scanner.AuditPlanSQLReq
	now := time.Now()
	auditPlanSqlMap := make(map[string]*scanner.AuditPlanSQLReq, 0)

	var sqlIdentity sqlIdentity
	var processedRawText, processedFingerPrint string
	for _, sql := range sqls {
		processedRawText = sql.RawText
		fingerPrint, err := pgParser.Normalize(sql.RawText)
		// 获取sql指纹失败时，保留SQL原语句
		if err != nil {
			processedFingerPrint = sql.RawText
		} else {
			processedFingerPrint = fingerPrint
		}

		sqlIdentity.Fingerprint = processedFingerPrint
		sqlIdentity.Schema = sql.Schema
		key := sqlIdentity.jsonString()

		if sqlReq, ok := auditPlanSqlMap[key]; ok {
			atoi, err := strconv.Atoi(sqlReq.Counter)
			if err != nil {
				return err
			}

			tMax := utils.MaxFloat64(*sqlReq.QueryTimeMax, sql.QueryTime)
			tAvg := utils.IncrementalAverageFloat64(*sqlReq.QueryTimeAvg, sql.QueryTime, atoi, 1)
			sqlReq.QueryTimeMax = &tMax
			sqlReq.QueryTimeAvg = &tAvg
			sqlReq.Counter = strconv.Itoa(atoi + 1)
			sqlReq.LastReceiveText = processedRawText
			sqlReq.LastReceiveTimestamp = now.Format(time.RFC3339)
			if sql.Endpoint != "" {
				sqlReq.Endpoints = append(sqlReq.Endpoints, sql.Endpoint)
			}
		} else {
			sqlReq := &scanner.AuditPlanSQLReq{
				Fingerprint:          processedFingerPrint,
				LastReceiveText:      processedRawText,
				LastReceiveTimestamp: now.Format(time.RFC3339),
				Counter:              "1",
				QueryTimeAvg:         &sql.QueryTime,
				QueryTimeMax:         &sql.QueryTime,
				FirstQueryAt:         sql.QueryAt,
				DBUser:               sql.DBUser,
				Schema:               sql.Schema,
			}
			if sql.Endpoint != "" {
				sqlReq.Endpoints = append(sqlReq.Endpoints, sql.Endpoint)
			}

			auditPlanSqlMap[key] = sqlReq
			sqlListReq = append(sqlListReq, sqlReq)
		}
	}

	return t.c.UploadReq(scanner.UploadSQL, t.auditPlanID, sqlListReq)
}
