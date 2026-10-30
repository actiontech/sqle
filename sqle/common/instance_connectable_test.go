package common

import (
	"testing"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

func TestCatalogDatabaseNameFromInstanceParams(t *testing.T) {
	tests := []struct {
		name   string
		params params.Params
		want   string
	}{
		{
			name: "AC-002 valued catalog",
			params: params.Params{
				&params.Param{Key: "database_name", Value: "testdb", Type: params.ParamTypeString},
				&params.Param{Key: "default_schema", Value: "public", Type: params.ParamTypeString},
			},
			want: "testdb",
		},
		{
			name: "AC-003 empty catalog",
			params: params.Params{
				&params.Param{Key: "database_name", Value: "", Type: params.ParamTypeString},
			},
			want: "",
		},
		{
			name: "AC-004 default_schema alone does not become database",
			params: params.Params{
				&params.Param{Key: "default_schema", Value: "public", Type: params.ParamTypeString},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &model.Instance{AdditionalParams: tt.params}
			got := catalogDatabaseNameFromInstance(inst)
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestSetParamValueStillRejectsUnknown(t *testing.T) {
	// AC-005: 未做全局 ignore —— 未知参数仍明确失败
	ps := params.Params{
		&params.Param{Key: "default_schema", Value: "public", Type: params.ParamTypeString},
		&params.Param{Key: "database_name", Value: "", Type: params.ParamTypeString},
	}
	err := ps.SetParamValue("databas_name", "x")
	if err == nil {
		t.Fatal("expected unknown param to fail")
	}
	if err.Error() != "param databas_name not found" {
		t.Fatalf("unexpected err: %v", err)
	}
}
