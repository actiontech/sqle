package dmslowlog

import (
	"bufio"
	"context"
	"fmt"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type DmSlowLog struct {
	l          *logrus.Entry
	c          *scanner.Client
	sqls       []scanners.SQL
	allSQL     []DmNode
	getAll     chan struct{}
	apName     string
	slowLogDir string // 同步达梦目录，如：/dm8/dmdbms/log
}

type Params struct {
	SlowLogDir             string // 同步达梦目录，如：/dm8/dmdbms/log
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
		slowLogDir: params.SlowLogDir,
		apName:     params.APName,
		l:          l,
		c:          c,
		getAll:     make(chan struct{}),
	}, nil
}

func (dm *DmSlowLog) Run(ctx context.Context) error {
	sqls, err := getSQLFromDir(dm.slowLogDir, scanners.LOGFileSuffix)
	if err != nil {
		return err
	}

	dm.allSQL = sqls
	close(dm.getAll)

	<-ctx.Done()
	return nil
}

func (dm *DmSlowLog) SQLs() <-chan scanners.SQL {
	sqlCh := make(chan scanners.SQL, 10240)

	go func() {
		<-dm.getAll
		for _, sql := range dm.allSQL {
			sqlCh <- scanners.SQL{
				Fingerprint: sql.Fingerprint,
				RawText:     sql.Text,
				Schema:      sql.Schema,
				QueryTime:   sql.QueryTime,
				QueryAt:     sql.QueryAt,
				DBUser:      sql.DBUser,
			}
		}
		close(sqlCh)
	}()
	return sqlCh
}

func (dm *DmSlowLog) Upload(ctx context.Context, sqls []scanners.SQL) error {
	dm.sqls = append(dm.sqls, sqls...)
	err := common.UploadForDmSlowLog(ctx, dm.sqls, dm.c, dm.apName)
	if err != nil {
		return err
	}
	return common.Audit(dm.c, dm.apName)
}

func getSQLFromDir(dir string, fileSuffix string) ([]DmNode, error) {
	dmNodes := make([]DmNode, 0)
	if fileSuffix == scanners.LOGFileSuffix {
		fileName, err := getLatestFile(dir, "dmsql_")
		if err != nil {
			return nil, err
		}
		dmsLogFileName := fmt.Sprintf("%s/%s", dir, fileName)
		file, err := tail.OpenFile(dmsLogFileName)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("close file=%s err:%s\n", file.Name(), err)
			}
		}(file)

		newScanner := bufio.NewScanner(file)
		for newScanner.Scan() {
			node := DmNode{}
			line := newScanner.Text()
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
					return nil, err
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
					return nil, err
				}
				node.QueryAt = queryAt
			}

			// user:SYSDBA trxid:
			userStart, userEnd := strings.Index(line, "user:")+5, strings.Index(line, "trxid:")
			// 执行SQL的用户
			node.DBUser = strings.TrimSpace(line[userStart:userEnd])

			// 解析sql和sql类型
			dmNode := parseSqlAndType(line)
			if dmNode != nil {
				node.Type, node.Text = dmNode.Type, dmNode.Text
				node.Fingerprint, err = util.Fingerprint(dmNode.Text, true)
				if err != nil {
					return nil, err
				}
				dmNodes = append(dmNodes, node)
			}
		}
	}
	return dmNodes, nil
}

func getLatestFile(dir, prefix string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("directiory=%s is empty", dir)
	}

	maxFileNumber := "0"
	maxFileIndex := -1
	for i, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if strings.HasPrefix(filename, prefix) {
			re := regexp.MustCompile(`\d+`)
			matches := re.FindAllString(strings.ReplaceAll(filename, "_", ""), -1)
			if len(matches) == 1 && strings.Compare(matches[0], maxFileNumber) > 0 {
				maxFileNumber = matches[0]
				maxFileIndex = i
			}
		}
	}

	if maxFileIndex == -1 {
		return "", fmt.Errorf("under directiory=%s does not exits dm slow log file", dir)
	}
	return files[maxFileIndex].Name(), nil
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
