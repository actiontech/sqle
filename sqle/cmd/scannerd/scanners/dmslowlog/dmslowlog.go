package dmslowlog

import (
	"context"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/hpcloud/tail"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type DmSlowLog struct {
	l           *logrus.Entry
	c           *scanner.Client
	sqls        []scanners.SQL
	sqlCh       chan scanners.SQL
	apName      string
	slowLogFile string // 同步达梦目录，如：/dm8/dmdbms/log
}

type Params struct {
	SlowLogFile            string // 同步达梦目录，如：/dm8/dmdbms/log
	APName                 string
	SyncDmSlowLogStartTime string // 同步达梦日志开始时间
	SkipErrorQuery         bool
	SkipErrorSqlFile       bool
	SkipAudit              bool
}

type DmNode struct {
	driverV2.Node
	Schema    string
	QueryTime float64   // 达梦慢日志执行时长
	QueryAt   time.Time // 达梦慢日志发生时间
	DBUser    string    // 执行SQL的用户
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*DmSlowLog, error) {
	return &DmSlowLog{
		slowLogFile: params.SlowLogFile,
		apName:      params.APName,
		sqls:        make([]scanners.SQL, 0),
		sqlCh:       make(chan scanners.SQL, 10),
		l:           l,
		c:           c,
	}, nil
}

func (dm *DmSlowLog) Run(ctx context.Context) error {
	tf, err := tail.TailFile(dm.slowLogFile, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Printf("open slow log file=%s error:%s\n", tf.Filename, err)
		return err
	}
	defer tf.Stop()

	for {
		select {
		case lines := <-tf.Lines:
			node := DmNode{}
			line := lines.Text
			// 正则表达式
			timestampRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}`)
			//sqlRegex := regexp.MustCompile(`([DDL]|[INS])* (.*);`)
			execTimeRegex := regexp.MustCompile(`EXECTIME: (\d+)(ms)*`)
			// 获取时间戳、SQL和执行时间
			timestamp := timestampRegex.FindStringSubmatch(line)
			//sql := sqlRegex.FindStringSubmatch(input)
			execTime := execTimeRegex.FindStringSubmatch(line)

			node.Schema = ""
			// 慢日志执行时长
			if len(execTime) == 3 {
				execTimeLong, err := strconv.ParseFloat(execTime[1], 64)
				if err != nil {
					log.Printf("parse string for float err:%s\n", err)
					continue
				}
				node.QueryTime = execTimeLong
				if execTimeLong == 0 {
					continue
				}
			}
			// 慢日志发生时间
			if len(timestamp) == 1 {
				queryAt, err := time.Parse("2006-01-02 15:04:05.999", timestamp[0])
				if err != nil {
					log.Printf("parse string for time err:%s\n", err)
					continue
				}
				node.QueryAt = queryAt
			}

			// user:SYSDBA trxid:
			userStart, userEnd := strings.Index(line, "user:"), strings.Index(line, "trxid:")
			if userStart == -1 || userEnd == -1 {
				continue
			}
			// 执行SQL的用户
			node.DBUser = strings.TrimSpace(line[userStart+5 : userEnd])

			// 解析sql和sql类型
			dmNode := parseSqlAndType(line)
			if dmNode != nil {
				node.Type, node.Text = dmNode.Type, dmNode.Text
				fingerprint, err := util.Fingerprint(dmNode.Text, true)
				if err != nil {
					log.Printf("parse sql finger err:%s\n", err)
					continue
				}
				node.Fingerprint = fingerprint
				scannerSql := scanners.SQL{
					Fingerprint: node.Fingerprint,
					RawText:     node.Text,
					Schema:      node.Schema,
					QueryTime:   node.QueryTime,
					QueryAt:     node.QueryAt,
					DBUser:      node.DBUser,
				}
				dm.sqls = append(dm.sqls, scannerSql)
				dm.sqlCh <- scannerSql
			}
		case <-ctx.Done():
			close(dm.sqlCh)
			return nil
		}
	}
}

func (dm *DmSlowLog) SQLs() <-chan scanners.SQL {
	return dm.sqlCh
}

func (dm *DmSlowLog) Upload(ctx context.Context, sqls []scanners.SQL) error {
	dm.sqls = append(dm.sqls, sqls...)
	err := common.UploadForDmSlowLog(ctx, dm.sqls, dm.c, dm.apName)
	if err != nil {
		return err
	}
	return common.Audit(dm.c, dm.apName)
}

func parseSqlAndType(line string) *DmNode {
	node := &DmNode{}
	isNeedInsertDb := true
	sql := ""
	sqlTypeSlice := []string{"[DDL]", "[ORA]", "[INS]", "[DEL]", "[UPD]", "[SEL]", "[UNKNOWN]"}
	for _, sqlType := range sqlTypeSlice {
		start, end := strings.Index(line, sqlType), strings.Index(line, "EXECTIME:")
		if start == -1 || end <= 0 {
			continue
		}
		sql = strings.TrimSpace(line[start+len(sqlType) : end])
		if len(sql) == 0 {
			isNeedInsertDb = false
			continue
		}
		switch sqlType {
		case "[DDL]", "[UNKNOWN]":
			node.Type = driverV2.SQLTypeDDL
		case "[ORA]", "[INS]", "[DEL]", "[UPD]", "[SEL]":
			node.Type = driverV2.SQLTypeDML
		default:
			node.Type = driverV2.SQLTypeDDL
		}
	}
	if isNeedInsertDb && len(sql) > 0 {
		node.Text = sql
		return node
	}
	return nil
}
