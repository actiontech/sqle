package state

import "github.com/guregu/null"

// PostgresFunction - Function/Stored Procedure that runs on the PostgreSQL server
type PostgresFunction struct {
	Oid             Oid
	DatabaseOid     Oid
	SchemaName      string      `json:"schema_name"`
	FunctionName    string      `json:"function_name"`
	Language        string      `json:"language"`
	Source          string      `json:"source"`
	SourceBin       null.String `json:"source_bin"`
	Config          []string    `json:"config"`
	Arguments       string      `json:"arguments"`
	Result          string      `json:"result"`
	Kind            string      `json:"kind"`
	SecurityDefiner bool        `json:"security_definer"`
	Leakproof       bool        `json:"leakproof"`
	Strict          bool        `json:"strict"`
	ReturnsSet      bool        `json:"returns_set"`
	Volatile        string      `json:"volatile"`
}

// PostgresFunctionStats - Statistics about a single PostgreSQL function
//
// Note that this will only be populated when "track_functions" is enabled.
type PostgresFunctionStats struct {
	Calls     int64   `json:"calls"`
	TotalTime float64 `json:"total_time"`
	SelfTime  float64 `json:"self_time"`
}

type PostgresFunctionStatsMap map[Oid]PostgresFunctionStats

type DiffedPostgresFunctionStats PostgresFunctionStats
type DiffedPostgresFunctionStatsMap map[Oid]DiffedPostgresFunctionStats

func (curr PostgresFunctionStats) DiffSince(prev PostgresFunctionStats) DiffedPostgresFunctionStats {
	return DiffedPostgresFunctionStats{
		Calls:     curr.Calls - prev.Calls,
		TotalTime: curr.TotalTime - prev.TotalTime,
		SelfTime:  curr.SelfTime - prev.SelfTime,
	}
}
