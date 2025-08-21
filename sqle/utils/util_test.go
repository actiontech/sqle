package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	type args struct {
		s             string
		prefix        string
		caseSensitive bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{s: "hello, world", prefix: "hel", caseSensitive: true}, true},
		{"", args{s: "hello, world", prefix: "HEL", caseSensitive: true}, false},
		{"", args{s: "hello, world", prefix: "hel", caseSensitive: false}, true},
		{"", args{s: "hello, world", prefix: "HEL", caseSensitive: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPrefix(tt.args.s, tt.args.prefix, tt.args.caseSensitive); got != tt.want {
				t.Errorf("HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	type args struct {
		s             string
		suffix        string
		caseSensitive bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{s: "hello, world", suffix: "rld", caseSensitive: true}, true},
		{"", args{s: "hello, world", suffix: "RLD", caseSensitive: true}, false},
		{"", args{s: "hello, world", suffix: "rld", caseSensitive: false}, true},
		{"", args{s: "hello, world", suffix: "RLD", caseSensitive: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.args.s, tt.args.suffix, tt.args.caseSensitive); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDuplicate(t *testing.T) {
	assert.Equal(t, []string{}, GetDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"2"}, GetDuplicate([]string{"1", "2", "2"}))
	assert.Equal(t, []string{"2", "3"}, GetDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestRemoveDuplicate(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, RemoveDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, RemoveDuplicate([]string{"1", "2", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, RemoveDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestRound(t *testing.T) {
	assert.Equal(t, float64(1), Round(1.11, 0))
	assert.Equal(t, float64(0), Round(1.111117, -2))
	assert.Equal(t, 1.1, Round(1.11, 1))
	assert.Equal(t, 1.11112, Round(1.111117, 5))
}

func TestSupplementalQuotationMarks(t *testing.T) {
	assert.Equal(t, "'asdf'", SupplementalQuotationMarks("'asdf'"))
	assert.Equal(t, "\"asdf\"", SupplementalQuotationMarks("\"asdf\""))
	assert.Equal(t, "`asdf`", SupplementalQuotationMarks("`asdf`"))
	assert.Equal(t, "", SupplementalQuotationMarks(""))
	assert.Equal(t, "`asdf`", SupplementalQuotationMarks("asdf"))
	assert.Equal(t, "`\"asdf`", SupplementalQuotationMarks("\"asdf"))
	assert.Equal(t, "`asdf\"`", SupplementalQuotationMarks("asdf\""))
	assert.Equal(t, "`'asdf`", SupplementalQuotationMarks("'asdf"))
	assert.Equal(t, "`asdf'`", SupplementalQuotationMarks("asdf'"))
	assert.Equal(t, "``asdf`", SupplementalQuotationMarks("`asdf"))
	assert.Equal(t, "`asdf``", SupplementalQuotationMarks("asdf`"))
	assert.Equal(t, "`\"asdf'`", SupplementalQuotationMarks("\"asdf'"))
	assert.Equal(t, "`\"asdf``", SupplementalQuotationMarks("\"asdf`"))
	assert.Equal(t, "`'asdf\"`", SupplementalQuotationMarks("'asdf\""))
	assert.Equal(t, "`'asdf``", SupplementalQuotationMarks("'asdf`"))
	assert.Equal(t, "``asdf\"`", SupplementalQuotationMarks("`asdf\""))
	assert.Equal(t, "``asdf'`", SupplementalQuotationMarks("`asdf'"))
	assert.Equal(t, "`s`", SupplementalQuotationMarks("s"))
}

func TestIsUpperAndLowerLetterMixed(t *testing.T) {
	type args struct {
		s    string
		want bool
	}
	tests := []args{
		{"isUPPER", true},
		{"ISupper", true},
		{"isUpper", true},
		{"___isUPPER", true},
		{"isUPPER__@@", true},
		{"isUPPER@!$and", true},
		{"process", false},
		{"___process", false},
		{"process!@#", false},
		{"process__@@cc", false},
		{"a", false},
		{"$", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := IsUpperAndLowerLetterMixed(tt.s); got != tt.want {
				t.Errorf("IsUpperAndLowerLetterMixed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEventSQL(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"with DELIMITER", args{sql: `-- 1. 修改分隔符为 $$
DELIMITER $$ 

-- 2. 编写包含多条语句的事件
CREATE EVENT my_multi_statement_event
ON SCHEDULE EVERY 1 DAY
DO
BEGIN
    -- 内部的第一条语句
    DELETE FROM old_logs WHERE log_date < NOW() - INTERVAL 1 YEAR;

    -- 内部的第二条语句
    INSERT INTO report_log (message) VALUES ('Old logs have been deleted.');
END $$

-- 3. 恢复分隔符为默认的分号
DELIMITER ;`}, true},
		{"test2", args{sql: "create event my_event on schedule every 10 second do update schema.table set mycol = mycol + 1;"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsEventSQL(tt.args.sql), "IsEventSQL(%v)", tt.args.sql)
		})
	}
}
