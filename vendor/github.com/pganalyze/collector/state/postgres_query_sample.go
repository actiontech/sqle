package state

import (
	"encoding/json"
	"time"

	"github.com/guregu/null"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	uuid "github.com/satori/go.uuid"
)

type ExplainPlanTrigger struct {
	Calls          *int64   `json:"Calls,omitempty"`
	ConstraintName *string  `json:"Constraint Name,omitempty"`
	Relation       *string  `json:"Relation,omitempty"`
	Time           *float64 `json:"Time,omitempty"`
	TriggerName    *string  `json:"Trigger Name,omitempty"`
}

type ExplainPlanJIT struct {
	Functions *int64              `json:"Functions,omitempty"`
	Options   *map[string]bool    `json:"Options,omitempty"`
	Timing    *map[string]float64 `json:"Timing,omitempty"`
}

type ExplainPlanContainer struct {
	ExecutionTime   *float64              `json:"Execution Time,omitempty"`
	JIT             *ExplainPlanJIT       `json:"JIT,omitempty"`
	Plan            json.RawMessage       `json:"Plan"`
	Planning        *map[string]int64     `json:"Planning,omitempty"`
	PlanningTime    *float64              `json:"Planning Time,omitempty"`
	QueryIdentifier *int64                `json:"Query Identifier,omitempty"`
	QueryText       string                `json:"Query Text,omitempty"`
	Settings        *map[string]string    `json:"Settings,omitempty"`
	Triggers        *[]ExplainPlanTrigger `json:"Triggers,omitempty"`
}

type PostgresQuerySample struct {
	OccurredAt time.Time
	Username   string
	Database   string
	Query      string
	Parameters []null.String

	LogLineUUID uuid.UUID

	RuntimeMs float64

	HasExplain        bool
	ExplainOutputText string
	ExplainOutputJSON *ExplainPlanContainer
	ExplainError      string
	ExplainFormat     pganalyze_collector.QuerySample_ExplainFormat
	ExplainSource     pganalyze_collector.QuerySample_ExplainSource
}
