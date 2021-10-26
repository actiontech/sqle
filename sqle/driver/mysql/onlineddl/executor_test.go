package onlineddl

import (
	"testing"

	_ "github.com/pingcap/tidb/types/parser_driver"
)

func Test_parseAlterTableOptions(t *testing.T) {
	tests := []struct {
		alter         string
		wantSchema    string
		wantTable     string
		wantAlterOpts string
		wantErr       bool
	}{
		{
			alter:         "alter table t1 add column i int, drop column d",
			wantSchema:    "",
			wantTable:     "t1",
			wantAlterOpts: "ADD COLUMN `i` INT,DROP COLUMN `d`",
			wantErr:       false,
		},

		{
			alter:         "alter table `db1`.`t1` add column i int, drop column d",
			wantSchema:    "db1",
			wantTable:     "t1",
			wantAlterOpts: "ADD COLUMN `i` INT,DROP COLUMN `d`",
			wantErr:       false,
		},

		{
			alter:         "alter table `db1`.`t1` add column i int",
			wantSchema:    "db1",
			wantTable:     "t1",
			wantAlterOpts: "ADD COLUMN `i` INT",
			wantErr:       false,
		},

		{
			alter:         "alter table t1 add column col4 varchar(2);",
			wantSchema:    "",
			wantTable:     "t1",
			wantAlterOpts: "ADD COLUMN `col4` VARCHAR(2)",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotSchema, gotTable, gotAlterOpts, err := parseAlterTableOptions(tt.alter)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAlterTableOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSchema != tt.wantSchema {
				t.Errorf("parseAlterTableOptions() gotSchema = %v, want %v", gotSchema, tt.wantSchema)
			}
			if gotTable != tt.wantTable {
				t.Errorf("parseAlterTableOptions() gotTable = %v, want %v", gotTable, tt.wantTable)
			}
			if gotAlterOpts != tt.wantAlterOpts {
				t.Errorf("parseAlterTableOptions() gotAlterOpts = %v, want %v", gotAlterOpts, tt.wantAlterOpts)
			}
		})
	}
}
