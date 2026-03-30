package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type SQLLineageAnalyzeReqV1 struct {
	SQL               string   `json:"sql" form:"sql" valid:"required"`
	InstanceType      string   `json:"instance_type" form:"instance_type" valid:"omitempty,oneof=MySQL"`
	DefaultSchema     string   `json:"default_schema" form:"default_schema"`
	ResultColumnNames []string `json:"result_columns" form:"result_columns"`
}

type SQLLineageAnalyzeResV1 struct {
	controller.BaseRes
	Data *SQLLineageAnalyzeResDataV1 `json:"data"`
}

type SQLLineageAnalyzeResDataV1 struct {
	Result  *SQLLineageAnalyzeResultV1 `json:"result"`
}

// SQLLineageAnalyzeResultV1 is the API-level analyze result.
// NOTE: Keep it decoupled from server/sql_lineage_analysis to allow CE build.
type SQLLineageAnalyzeResultV1 struct {
	Title         string                 `json:"title"`
	OriginalSQL   string                 `json:"original_sql"`
	Tables        []SQLLineageTableRefV1 `json:"tables"`
	SourceColumns []SQLLineageColumnRefV1 `json:"source_columns"`
	ResultColumns []SQLLineageResultColumnV1 `json:"result_columns"`
	Nodes         []SQLLineageNodeV1     `json:"nodes"`
	Edges         []SQLLineageEdgeV1     `json:"edges"`
	Warnings      []string               `json:"warnings"`
}

type SQLLineageTableRefV1 struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Alias  string `json:"alias"`
}

type SQLLineageColumnRefV1 struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Column string `json:"column"`
}

type SQLLineageResultColumnV1 struct {
	Name       string                 `json:"name"`
	Expression string                 `json:"expression"`
	Sources    []SQLLineageColumnRefV1 `json:"sources"`
}

type SQLLineageNodeV1 struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Schema string `json:"schema,omitempty"`
	Table  string `json:"table,omitempty"`
	Column string `json:"column,omitempty"`
	Expr   string `json:"expr,omitempty"`
}

type SQLLineageEdgeV1 struct {
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
	Type   string `json:"type"`
}

// SQLLineageAnalyze
// @Summary SQL列级血缘分析
// @Description Analyze SQL and return column-level lineage
// @Id sqlLineageAnalyzeV1
// @Tags sql_analysis
// @Security ApiKeyAuth
// @Param req body v1.SQLLineageAnalyzeReqV1 true "sql lineage analyze request"
// @Success 200 {object} v1.SQLLineageAnalyzeResV1
// @Router /v1/sql_lineage_analysis [post]
func SQLLineageAnalyze(c echo.Context) error {
	return sqlLineageAnalyze(c)
}

