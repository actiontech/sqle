//go:build enterprise
// +build enterprise

package tbase_audit_log

import (
	"database/sql"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

type testTbaseLog struct {
	TxStartTime        string
	SqlType            string
	User               string
	Schema             string
	ClientHostWithPort string
	ConnectionTime     string
	Duration           float64
	SQLText            string
}

type info struct {
	id           int
	exceptResult *testTbaseLog
	message      string
	parseError   string
	logContent   string
}

func Test_ParseSql(t *testing.T) {
	infos := []info{
		{
			id:           0,
			logContent:   `2024-06-25 11:00:00.715 CST,"tbase","postgres",3067712,coord(3067712,6014334),"158.219.101.86:60174",667a32b0.2ecf40,coord(3067712,6014334),coord(0,0),2,"authentication",2024-06-25 11:00:00 CST,23/6014334,0,LOG,00000,"connection authorized: user=tbase database=postgres",,,,,,,,,""`,
			exceptResult: nil,
			message:      "connection log",
		},
		{
			id:           1,
			logContent:   `2024-06-25 11:00:00.718 CST,"tbase","postgres",3067712,coord(3067712,0),"158.219.101.86:60174",667a32b0.2ecf40,coord(3067712,0),coord(0,0),3,"idle",2024-06-25 11:00:00 CST,,0,LOG,00000,"disconnection: session time: 0:00:00.004 user=tbase database=postgres host=158.219.101.86 port=60174",,,,,,,,,""`,
			exceptResult: nil,
			message:      "disconnection log",
		},
		{
			id:           2,
			logContent:   `Stack trace:`,
			exceptResult: nil,
			message:      "database error log",
			parseError:   "record is not a standard log output, record: [Stack trace:]",
		},
		{
			id:           3,
			logContent:   `1    0xeae29e postgres errstart + 0x52e`,
			exceptResult: nil,
			message:      "database error log",
			parseError:   "record is not a standard log output, record: [1    0xeae29e postgres errstart + 0x52e]",
		},
		{
			id:           4,
			logContent:   `Use addr2line to get pretty function name and line, e.g., addr2line -e path_to_postgres symbol_address -f`,
			exceptResult: nil,
			message:      "database error log",
			parseError:   "record is not a standard log output, record: [Use addr2line to get pretty function name and line  e.g.  addr2line -e path_to_postgres symbol_address -f]",
		},
		{
			id:         5,
			logContent: `2024-06-25 11:11:25.521 CST,"script_user","dqsp",3117387,coord(3117387,411293),"158.220.194.22:16388",667a355c.2f914b,coord(3117387,411293),coord(0,0),3,"SELECT",2024-06-25 11:11:24 CST,43/411293,0,LOG,00000,"duration: 1086.036 ms  execute <unnamed>: select * from test where card_no=$1 and tran_dt>=$2 and tran_dt<=$3","parameters: $1 = '12', $2 = '20191027', $3 = '20231231'",,,,,,,,"PostgreSQL JDBC Driver"`,
			exceptResult: &testTbaseLog{
				TxStartTime:        "2024-06-25 11:11:25.521 +0000 UTC",
				SqlType:            "SELECT",
				User:               "script_user",
				Schema:             "dqsp",
				ClientHostWithPort: "158.220.194.22:16388",
				ConnectionTime:     "2024-06-25 11:11:24 +0000 UTC",
				Duration:           1086.036,
				SQLText:            "select * from test where card_no='12' and tran_dt>='20191027' and tran_dt<='20231231'",
			},
			message: "slow sql with all info",
		},
		{
			id:         6,
			logContent: `2024-06-25 11:11:25.521 CST,"script_user","dqsp",3117387,,"158.220.194.22:16388",667a355c.2f914b,,,3,"SELECT",2024-06-25 11:11:24 CST,43/411293,0,LOG,00000,"duration: 1086.036 ms  execute <unnamed>: select * from test where card_no=$1 and tran_dt>=$2 and tran_dt<=$3","parameters: $1 = '12', $2 = '20191027', $3 = '20231231'",,,,,,,,"PostgreSQL JDBC Driver"`,
			exceptResult: &testTbaseLog{
				TxStartTime:        "2024-06-25 11:11:25.521 +0000 UTC",
				SqlType:            "SELECT",
				User:               "script_user",
				Schema:             "dqsp",
				ClientHostWithPort: "158.220.194.22:16388",
				ConnectionTime:     "2024-06-25 11:11:24 +0000 UTC",
				Duration:           1086.036,
				SQLText:            "select * from test where card_no='12' and tran_dt>='20191027' and tran_dt<='20231231'",
			},
			message: "slow sql without coord value",
		},
		{
			id:         7,
			logContent: `2024-06-25 11:11:25.521 CST,"script_user","dqsp",3117387,coord(3117387,411293),"158.220.194.22:16388",667a355c.2f914b,coord(3117387,411293),coord(0,0),3,"SELECT",2024-06-25 11:11:24 CST,43/411293,0,LOG,00000,"duration: 1086.036 ms  execute <unnamed>: select * from test",,,,,,,,,"PostgreSQL JDBC Driver"`,
			exceptResult: &testTbaseLog{
				TxStartTime:        "2024-06-25 11:11:25.521 +0000 UTC",
				SqlType:            "SELECT",
				User:               "script_user",
				Schema:             "dqsp",
				ClientHostWithPort: "158.220.194.22:16388",
				ConnectionTime:     "2024-06-25 11:11:24 +0000 UTC",
				Duration:           1086.036,
				SQLText:            "select * from test",
			},
			message: "slow sql without params",
		},
		{
			id: 8,
			logContent: `2024-06-25 11:11:25.521 CST,"script_user","dqsp",3117387,coord(3117387,411293),"158.220.194.22:16388",667a355c.2f914b,coord(3117387,411293),coord(0,0),3,"SELECT",2024-06-25 11:11:24 CST,43/411293,0,LOG,00000,"duration: 1086.036 ms  execute <unnamed>: select * 
			from test where id=1",,,,,,,,,"PostgreSQL JDBC Driver"`,
			exceptResult: &testTbaseLog{
				TxStartTime:        "2024-06-25 11:11:25.521 +0000 UTC",
				SqlType:            "SELECT",
				User:               "script_user",
				Schema:             "dqsp",
				ClientHostWithPort: "158.220.194.22:16388",
				ConnectionTime:     "2024-06-25 11:11:24 +0000 UTC",
				Duration:           1086.036,
				SQLText: `select * 
			from test where id=1`,
			},
			message: "slow sql with line break",
		},
	}

	parser := NewTbaseParser()

	for i, info := range infos {
		reader := strings.NewReader(info.logContent)
		csvReader := csv.NewReader(reader)
		record, err := csvReader.Read()
		if err != nil {
			t.Error(err)
		}

		tlog, err := parser.parseRecord(record)
		if err != nil {
			assert.EqualError(t, err, info.parseError)
		}
		if tlog == nil {
			if info.exceptResult != nil {
				t.Errorf("Expected value does not match actual value, infos id: %v", i)
			}
			continue
		}
		result := compareStruct(info.exceptResult, tlog)
		assert.Equalf(t, result, true, "Incorrect comparison results, infos id: %v", i)
	}
}

func compareStruct(testLog *testTbaseLog, tlog *TbaseLog) bool {
	if testLog.TxStartTime == tlog.TxStartTime.String() && testLog.Schema == tlog.Schema && testLog.User == tlog.User && testLog.ClientHostWithPort == tlog.ClientHostWithPort && testLog.Duration == tlog.Duration && testLog.SQLText == tlog.SQLText {
		return true
	}
	return false
}

type testSort struct {
	id            int
	testFileInfos []FileInfo
	exceptResult  []string
}

func Test_SortFiles(t *testing.T) {
	datas := []testSort{
		{
			id: 1,
			testFileInfos: []FileInfo{
				{
					Name:    "/test/Postgresql-0801-1.csv",
					ModTime: time.Date(2023, time.August, 1, 1, 20, 1, 1, time.Local), // 2023-08-01 01:20:01
				},
				{
					Name:    "/test/Postgresql-0801-2.csv",
					ModTime: time.Date(2023, time.August, 1, 3, 20, 1, 1, time.Local), // 2023-08-01 03:20:01
				},
				{
					Name:    "/test/Postgresql-0801-3.csv",
					ModTime: time.Date(2023, time.August, 1, 5, 20, 1, 1, time.Local), // 2023-08-01 05:20:01
				},
			},
			exceptResult: []string{"/test/Postgresql-0801-3.csv", "/test/Postgresql-0801-2.csv", "/test/Postgresql-0801-1.csv"},
		},
		{
			id: 2,
			testFileInfos: []FileInfo{
				{
					Name:    "/test/Postgresql-0801-3.csv",
					ModTime: time.Date(2023, time.August, 1, 5, 20, 1, 1, time.Local), // 2023-08-01 05:20:01
				},
				{
					Name:    "/test/Postgresql-0801-2.csv",
					ModTime: time.Date(2023, time.August, 1, 3, 20, 1, 1, time.Local), // 2023-08-01 03:20:01
				},
				{
					Name:    "/test/Postgresql-0801-1.csv",
					ModTime: time.Date(2023, time.August, 1, 1, 20, 1, 1, time.Local), // 2023-08-01 01:20:01
				},
			},
			exceptResult: []string{"/test/Postgresql-0801-3.csv", "/test/Postgresql-0801-2.csv", "/test/Postgresql-0801-1.csv"},
		},
		{
			id: 3,
			testFileInfos: []FileInfo{
				{
					Name:    "/test/Postgresql-0801-2.csv",
					ModTime: time.Date(2023, time.August, 1, 3, 20, 1, 1, time.Local), // 2023-08-01 03:20:01
				},
				{
					Name:    "/test/Postgresql-0801-3.csv",
					ModTime: time.Date(2023, time.August, 1, 5, 20, 1, 1, time.Local), // 2023-08-01 05:20:01
				},
				{
					Name:    "/test/Postgresql-0801-1.csv",
					ModTime: time.Date(2023, time.August, 1, 1, 20, 1, 1, time.Local), // 2023-08-01 01:20:01
				},
			},
			exceptResult: []string{"/test/Postgresql-0801-3.csv", "/test/Postgresql-0801-2.csv", "/test/Postgresql-0801-1.csv"},
		},
		{
			id: 4,
			testFileInfos: []FileInfo{
				{
					Name:    "/test/Postgresql-0801-2.csv",
					ModTime: time.Date(2023, time.August, 1, 3, 20, 1, 1, time.Local), // 2023-08-01 03:20:01
				},
				{
					Name:    "/test/Postgresql-0801-4.csv",
					ModTime: time.Date(2023, time.August, 1, 7, 20, 1, 1, time.Local), // 2023-08-01 07:20:01
				},
				{
					Name:    "/test/Postgresql-0801-3.csv",
					ModTime: time.Date(2023, time.August, 1, 5, 20, 1, 1, time.Local), // 2023-08-01 05:20:01
				},
				{
					Name:    "/test/Postgresql-0801-1.csv",
					ModTime: time.Date(2023, time.August, 1, 1, 20, 1, 1, time.Local), // 2023-08-01 01:20:01
				},
			},
			exceptResult: []string{"/test/Postgresql-0801-4.csv", "/test/Postgresql-0801-3.csv", "/test/Postgresql-0801-2.csv", "/test/Postgresql-0801-1.csv"},
		},
	}
	for _, data := range datas {
		sortedFiles := sortFilesByModTime(data.testFileInfos)
		for i, fileInfo := range sortedFiles {
			if fileInfo.Name != data.exceptResult[i] {
				t.Errorf("the files are not sorted as expected, id: %v", data.id)
			}
		}
	}
}

type testSqlWithParams struct {
	id           int
	sql          string
	params       []null.String
	message      string
	exceptResult string
}

func Test_ReplaceParams(t *testing.T) {
	datas := []testSqlWithParams{
		{
			id:           1,
			message:      "sql without params",
			sql:          "select * from test",
			params:       []null.String{},
			exceptResult: "select * from test",
		},
		{
			id:      2,
			message: "sql with params",
			sql:     "select * from test where name=$1 and data > $2",
			params: []null.String{
				{
					NullString: sql.NullString{
						Valid:  true,
						String: "abc",
					},
				},
				{
					NullString: sql.NullString{
						Valid:  true,
						String: "2019",
					},
				},
			},
			exceptResult: "select * from test where name='abc' and data > '2019'",
		},
		{
			id:      3,
			message: "sql with little params",
			sql:     "select * from test where name=$1 and data>$2",
			params: []null.String{
				{
					NullString: sql.NullString{
						Valid:  true,
						String: "abc",
					},
				},
			},
			exceptResult: "select * from test where name='abc' and data>$2",
		},
		{
			id:           4,
			message:      "sql with more params",
			sql:          "select * from test where name=$1",
			params:       []null.String{
				{
					NullString: sql.NullString{
						Valid:  true,
						String: "abc",
					},
				},
				{
					NullString: sql.NullString{
						Valid:  true,
						String: "2019",
					},
				},
			},
			exceptResult: "select * from test where name='abc'",
		},
	}
	for _, data := range datas {
		sql := replaceParams(data.sql, data.params)
		assert.Equalf(t, sql, data.exceptResult, "sql not match except result, id: %v", data.id)
	}
}
