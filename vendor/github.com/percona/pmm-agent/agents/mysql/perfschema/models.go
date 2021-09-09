// pmm-agent
// Copyright 2019 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package perfschema

//go:generate reform

// eventsStatementsSummaryByDigest represents a row in performance_schema.events_statements_summary_by_digest table.
//reform:performance_schema.events_statements_summary_by_digest
type eventsStatementsSummaryByDigest struct {
	SchemaName              *string `reform:"SCHEMA_NAME"`
	Digest                  *string `reform:"DIGEST"`      // MD5 of DigestText
	DigestText              *string `reform:"DIGEST_TEXT"` // query without values
	CountStar               uint64  `reform:"COUNT_STAR"`
	SumTimerWait            uint64  `reform:"SUM_TIMER_WAIT"`
	MinTimerWait            uint64  `reform:"MIN_TIMER_WAIT"`
	AvgTimerWait            uint64  `reform:"AVG_TIMER_WAIT"`
	MaxTimerWait            uint64  `reform:"MAX_TIMER_WAIT"`
	SumLockTime             uint64  `reform:"SUM_LOCK_TIME"`
	SumErrors               uint64  `reform:"SUM_ERRORS"`
	SumWarnings             uint64  `reform:"SUM_WARNINGS"`
	SumRowsAffected         uint64  `reform:"SUM_ROWS_AFFECTED"`
	SumRowsSent             uint64  `reform:"SUM_ROWS_SENT"`
	SumRowsExamined         uint64  `reform:"SUM_ROWS_EXAMINED"`
	SumCreatedTmpDiskTables uint64  `reform:"SUM_CREATED_TMP_DISK_TABLES"`
	SumCreatedTmpTables     uint64  `reform:"SUM_CREATED_TMP_TABLES"`
	SumSelectFullJoin       uint64  `reform:"SUM_SELECT_FULL_JOIN"`
	SumSelectFullRangeJoin  uint64  `reform:"SUM_SELECT_FULL_RANGE_JOIN"`
	SumSelectRange          uint64  `reform:"SUM_SELECT_RANGE"`
	SumSelectRangeCheck     uint64  `reform:"SUM_SELECT_RANGE_CHECK"`
	SumSelectScan           uint64  `reform:"SUM_SELECT_SCAN"`
	SumSortMergePasses      uint64  `reform:"SUM_SORT_MERGE_PASSES"`
	SumSortRange            uint64  `reform:"SUM_SORT_RANGE"`
	SumSortRows             uint64  `reform:"SUM_SORT_ROWS"`
	SumSortScan             uint64  `reform:"SUM_SORT_SCAN"`
	SumNoIndexUsed          uint64  `reform:"SUM_NO_INDEX_USED"`
	SumNoGoodIndexUsed      uint64  `reform:"SUM_NO_GOOD_INDEX_USED"`
	// FirstSeen               time.Time `reform:"FIRST_SEEN"`
	// LastSeen                time.Time `reform:"LAST_SEEN"`
}

// eventsStatementsSummaryByDigestExamples represents a row in
// performance_schema.events_statements_summary_by_digest table for MySQL 8.0 examples.
//reform:performance_schema.events_statements_summary_by_digest
type eventsStatementsSummaryByDigestExamples struct {
	SQLText       *string `reform:"QUERY_SAMPLE_TEXT"`
	Digest        *string `reform:"DIGEST"`
	CurrentSchema *string `reform:"SCHEMA_NAME"`
}

// eventsStatementsHistory represents a row in performance_schema.events_statements_history table.
//reform:performance_schema.events_statements_history
type eventsStatementsHistory struct {
	// ThreadID   int64   `reform:"THREAD_ID"`
	// EventID    int64   `reform:"EVENT_ID"`
	// EndEventID *int64  `reform:"END_EVENT_ID"`
	// EventName  string  `reform:"EVENT_NAME"`
	// Source     *string `reform:"SOURCE"`
	// TimerStart *int64  `reform:"TIMER_START"`
	// TimerEnd   *int64  `reform:"TIMER_END"`
	// TimerWait  *int64  `reform:"TIMER_WAIT"`
	// LockTime   int64   `reform:"LOCK_TIME"`
	SQLText *string `reform:"SQL_TEXT"`
	Digest  *string `reform:"DIGEST"`
	// DigestText    *string `reform:"DIGEST_TEXT"`
	CurrentSchema *string `reform:"CURRENT_SCHEMA"`
	// ObjectType           *string `reform:"OBJECT_TYPE"`
	// ObjectSchema         *string `reform:"OBJECT_SCHEMA"`
	// ObjectName           *string `reform:"OBJECT_NAME"`
	// ObjectInstanceBegin  *int64  `reform:"OBJECT_INSTANCE_BEGIN"`
	// MySQLErrno           *int32  `reform:"MYSQL_ERRNO"`
	// ReturnedSqlstate     *string `reform:"RETURNED_SQLSTATE"`
	// MessageText          *string `reform:"MESSAGE_TEXT"`
	// Errors               int64   `reform:"ERRORS"`
	// Warnings             int64   `reform:"WARNINGS"`
	// RowsAffected         int64   `reform:"ROWS_AFFECTED"`
	// RowsSent             int64   `reform:"ROWS_SENT"`
	// RowsExamined         int64   `reform:"ROWS_EXAMINED"`
	// CreatedTmpDiskTables int64   `reform:"CREATED_TMP_DISK_TABLES"`
	// CreatedTmpTables     int64   `reform:"CREATED_TMP_TABLES"`
	// SelectFullJoin       int64   `reform:"SELECT_FULL_JOIN"`
	// SelectFullRangeJoin  int64   `reform:"SELECT_FULL_RANGE_JOIN"`
	// SelectRange          int64   `reform:"SELECT_RANGE"`
	// SelectRangeCheck     int64   `reform:"SELECT_RANGE_CHECK"`
	// SelectScan           int64   `reform:"SELECT_SCAN"`
	// SortMergePasses      int64   `reform:"SORT_MERGE_PASSES"`
	// SortRange            int64   `reform:"SORT_RANGE"`
	// SortRows             int64   `reform:"SORT_ROWS"`
	// SortScan             int64   `reform:"SORT_SCAN"`
	// MoIndexUsed          int64   `reform:"NO_INDEX_USED"`
	// MoGoodIndexUsed      int64   `reform:"NO_GOOD_INDEX_USED"`
}

// setupConsumers represents a row in performance_schema.setup_consumers table.
//reform:performance_schema.setup_consumers
type setupConsumers struct {
	Name    string `reform:"NAME"`
	Enabled string `reform:"ENABLED"`
}

// setupInstruments represents a row in performance_schema.setup_instruments table.
//reform:performance_schema.setup_instruments
type setupInstruments struct {
	Name    string  `reform:"NAME"`
	Enabled string  `reform:"ENABLED"`
	Timed   *string `reform:"TIMED"` // nullable in 8.0
}
