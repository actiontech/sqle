package querysample

import (
	"encoding/json"
	"strings"

	"github.com/pganalyze/collector/state"
	"github.com/pganalyze/collector/util"
)

type planNodeIncrementalSort struct {
	GroupCount      *int64            `json:"Group Count,omitempty"`
	SortMethodsUsed *[]string         `json:"Sort Methods Used,omitempty"` // does not need normalize (fixed set of possible values)
	SortSpaceMemory *map[string]int64 `json:"Sort Space Memory,omitempty"`
	SortSpaceDisk   *map[string]int64 `json:"Sort Space Disk,omitempty"`
}

type planNodeGroupingSet struct {
	GroupKeys *[][]string `json:"Group Keys,omitempty"` // needs normalize (value comes from deparse_expression)
	HashKeys  *[][]string `json:"Hash Keys,omitempty"`  // needs normalize (value comes from deparse_expression)
	SortKey   *[]string   `json:"Sort Key,omitempty"`   // needs normalize (value comes from deparse_expression)
}

type planNode struct {
	ActualLoops                 *int64                   `json:"Actual Loops,omitempty"`
	ActualRows                  *int64                   `json:"Actual Rows,omitempty"`
	ActualStartupTime           *float64                 `json:"Actual Startup Time,omitempty"`
	ActualTotalTime             *float64                 `json:"Actual Total Time,omitempty"`
	Alias                       *string                  `json:"Alias,omitempty"` // does not need normalize (table alias as specified in the query)
	AsyncCapable                *bool                    `json:"Async Capable,omitempty"`
	CacheEvictions              *int64                   `json:"Cache Evictions,omitempty"`
	CacheHits                   *int64                   `json:"Cache Hits,omitempty"`
	CacheKey                    *string                  `json:"Cache Key,omitempty"` // needs normalize (value comes from deparse_expression)
	CacheMisses                 *int64                   `json:"Cache Misses,omitempty"`
	CacheMode                   *string                  `json:"Cache Mode,omitempty"` // does not need normalize (fixed set of possible values)
	CacheOverflows              *int64                   `json:"Cache Overflows,omitempty"`
	CTEName                     *string                  `json:"CTE Name,omitempty"`                 // does not need normalize (name of a CTE in the query)
	Command                     *string                  `json:"Command,omitempty"`                  // does not need normalize (fixed set of possible values)
	ConflictArbiterIndexes      *[]string                `json:"Conflict Arbiter Indexes,omitempty"` // does not need normalize (name of an index)
	ConflictFilter              *string                  `json:"Conflict Filter,omitempty"`          // needs normalize (value comes from show_upper_qual)
	ConflictResolution          *string                  `json:"Conflict Resolution,omitempty"`      // does not need normalize (fixed set of possible values)
	ConflictingTuples           *float64                 `json:"Conflicting Tuples,omitempty"`
	CustomPlanProvider          *string                  `json:"Custom Plan Provider,omitempty"` // does not need normalize (name of a custom plan provider implemented by an extension)
	DiskUsage                   *int64                   `json:"Disk Usage,omitempty"`
	ExactHeapBlocks             *int64                   `json:"Exact Heap Blocks,omitempty"`
	Filter                      *string                  `json:"Filter,omitempty"`           // needs normalize (value comes from show_scan_qual)
	FunctionCall                *string                  `json:"Function Call,omitempty"`    // needs normalize (value comes from show_expression)
	FunctionName                *string                  `json:"Function Name,omitempty"`    // does not neet normalize (name of a function)
	FullSortGroups              *planNodeIncrementalSort `json:"Full-sort Groups,omitempty"` // does not need nested normalize (see struct)
	GroupKey                    *[]string                `json:"Group Key,omitempty"`        // needs normalize (value comes from deparse_expression)
	GroupingSets                *[]planNodeGroupingSet   `json:"Grouping Sets,omitempty"`    // needs nested normalize (see struct)
	HashAggBatches              *int64                   `json:"HashAgg Batches,omitempty"`
	HashBatches                 *int64                   `json:"Hash Batches,omitempty"`
	HashBuckets                 *int64                   `json:"Hash Buckets,omitempty"`
	HashCond                    *string                  `json:"Hash Cond,omitempty"` // needs normalize (value comes from show_upper_qual)
	HeapFetches                 *int64                   `json:"Heap Fetches,omitempty"`
	IOReadTime                  *float64                 `json:"I/O Read Time,omitempty"`
	IOWriteTime                 *float64                 `json:"I/O Write Time,omitempty"`
	IndexCond                   *string                  `json:"Index Cond,omitempty"` // needs normalize (value comes from show_scan_qual)
	IndexName                   *string                  `json:"Index Name,omitempty"` // does not need normalize (name of an index)
	InnerUnique                 *bool                    `json:"Inner Unique,omitempty"`
	JoinFilter                  *string                  `json:"Join Filter,omitempty"` // needs normalize (value comes from show_upper_qual)
	JoinType                    *string                  `json:"Join Type,omitempty"`   // does not need normalize (fixed set of possible values)
	LocalDirtiedBlocks          *int64                   `json:"Local Dirtied Blocks,omitempty"`
	LocalHitBlocks              *int64                   `json:"Local Hit Blocks,omitempty"`
	LocalReadBlocks             *int64                   `json:"Local Read Blocks,omitempty"`
	LocalWrittenBlocks          *int64                   `json:"Local Written Blocks,omitempty"`
	LossyHeapBlocks             *int64                   `json:"Lossy Heap Blocks,omitempty"`
	MergeCond                   *string                  `json:"Merge Cond,omitempty"`      // needs normalize (value comes from show_upper_qual)
	NodeType                    *string                  `json:"Node Type,omitempty"`       // does not need normalize (fixed set of possible values)
	OneTimeFilter               *string                  `json:"One-Time Filter,omitempty"` // needs normalize (value comes from show_upper_qual)
	Operation                   *string                  `json:"Operation,omitempty"`       // does not need normalize (fixed set of possible values)
	OrderBy                     *string                  `json:"Order By,omitempty"`        // needs normalize (value comes from show_scan_qual)
	OriginalHashBatches         *int64                   `json:"Original Hash Batches,omitempty"`
	OriginalHashBuckets         *int64                   `json:"Original Hash Buckets,omitempty"`
	Output                      *[]string                `json:"Output,omitempty"` // needs normalize (value comes from deparse_expression)
	ParallelAware               *bool                    `json:"Parallel Aware,omitempty"`
	ParamsEvaluated             *[]string                `json:"Params Evaluated,omitempty"`    // does not need normalize (parameters are constructed explicitly as $n)
	ParentRelationship          *string                  `json:"Parent Relationship,omitempty"` // does not need normalize (fixed set of possible values)
	PartialMode                 *string                  `json:"Partial Mode,omitempty"`        // does not need normalize (fixed set of possible values)
	PeakMemoryUsage             *int64                   `json:"Peak Memory Usage,omitempty"`
	PlanRows                    *int64                   `json:"Plan Rows,omitempty"`
	PlanWidth                   *int64                   `json:"Plan Width,omitempty"`
	PlannedPartitions           *int64                   `json:"Planned Partitions,omitempty"`
	Plans                       []planNode               `json:"Plans,omitempty"`
	PreSortedGroups             *planNodeIncrementalSort `json:"Pre-sorted Groups,omitempty"` // does not need nested normalize (see struct)
	PresortedKey                *[]string                `json:"Presorted Key,omitempty"`     // needs normalize (value comes from deparse_expression)
	RecheckCond                 *string                  `json:"Recheck Cond,omitempty"`      // needs normalize (value comes from show_scan_qual)
	RelationName                *string                  `json:"Relation Name,omitempty"`     // does not need normalize (name of a table/view)
	RepeatableSeed              *string                  `json:"Repeatable Seed,omitempty"`   // needs normalize (value comes from deparse_expression)
	RowsRemovedByConflictFilter *int64                   `json:"Rows Removed by Conflict Filter,omitempty"`
	RowsRemovedByFilter         *int64                   `json:"Rows Removed by Filter,omitempty"`
	RowsRemovedByIndexRecheck   *int64                   `json:"Rows Removed by Index Recheck,omitempty"`
	RowsRemovedByJoinFilter     *int64                   `json:"Rows Removed by Join Filter,omitempty"`
	SamplingMethod              *string                  `json:"Sampling Method,omitempty"`     // does not need normalize (name of a function)
	SamplingParameters          *[]string                `json:"Sampling Parameters,omitempty"` // needs normalize (value comes from deparse_expression)
	ScanDirection               *string                  `json:"Scan Direction,omitempty"`      // does not need normalize (fixed set of possible values)
	Schema                      *string                  `json:"Schema,omitempty"`              // does not need normalize (name of a schema)
	SharedDirtiedBlocks         *int64                   `json:"Shared Dirtied Blocks,omitempty"`
	SharedHitBlocks             *int64                   `json:"Shared Hit Blocks,omitempty"`
	SharedReadBlocks            *int64                   `json:"Shared Read Blocks,omitempty"`
	SharedWrittenBlocks         *int64                   `json:"Shared Written Blocks,omitempty"`
	SingleCopy                  *bool                    `json:"Single Copy,omitempty"`
	SortKey                     *[]string                `json:"Sort Key,omitempty"`        // needs normalize (value comes from deparse_expression)
	SortMethod                  *string                  `json:"Sort Method,omitempty"`     // does not need normalize (fixed set of possible values)
	SortSpaceType               *string                  `json:"Sort Space Type,omitempty"` // does not need normalize (fixed set of possible values)
	SortSpaceUsed               *int64                   `json:"Sort Space Used,omitempty"`
	StartupCost                 *float64                 `json:"Startup Cost,omitempty"`
	Strategy                    *string                  `json:"Strategy,omitempty"`     // does not need normalize (fixed set of possible values)
	SubplanName                 *string                  `json:"Subplan Name,omitempty"` // does not need normalize (generated name by planner)
	SubplansRemoved             *int64                   `json:"Subplans Removed,omitempty"`
	TableFunctionCall           *string                  `json:"Table Function Call,omitempty"` // needs normalize (value comes from show_expression)
	TableFunctionName           *string                  `json:"Table Function Name,omitempty"` // does not need normalize (name of a function)
	TempReadBlocks              *int64                   `json:"Temp Read Blocks,omitempty"`
	TempWrittenBlocks           *int64                   `json:"Temp Written Blocks,omitempty"`
	TIDCond                     *string                  `json:"TID Cond,omitempty"`        // needs normalize (value comes from show_scan_qual)
	TuplestoreName              *string                  `json:"Tuplestore Name,omitempty"` // does not need normalize (name of an emphemeral named relation, aka ENR) -- note that there is no test exercising this, since its not clear how this part of explain.c can get called
	TotalCost                   *float64                 `json:"Total Cost,omitempty"`
	TuplesInserted              *float64                 `json:"Tuples Inserted,omitempty"`
	WALBytes                    *uint64                  `json:"WAL Bytes,omitempty"` // Note this is the only *unsigned* integer field currently
	WALFPI                      *int64                   `json:"WAL FPI,omitempty"`
	WALRecords                  *int64                   `json:"WAL Records,omitempty"`
	WorkerNumber                *int64                   `json:"Worker Number,omitempty"`
	Workers                     *[]planNode              `json:"Workers,omitempty"`
	WorkersLaunched             *int64                   `json:"Workers Launched,omitempty"`
	WorkersPlanned              *int64                   `json:"Workers Planned,omitempty"`
}

func normalizeExpr(expr *string) {
	if expr == nil {
		return
	}
	res := util.NormalizeQuery("SELECT "+*expr, "unparsable", -1)
	if res == util.QueryTextUnparsable {
		*expr = util.QueryTextUnparsable
		return
	}

	*expr = strings.TrimPrefix(res, "SELECT ")
}

func normalizeSortKey(expr *string) {
	if expr == nil {
		return
	}
	res := util.NormalizeQuery("SELECT ORDER BY "+*expr, "unparsable", -1)
	if res == util.QueryTextUnparsable {
		*expr = util.QueryTextUnparsable
		return
	}

	*expr = strings.TrimPrefix(res, "SELECT ORDER BY ")
}

func normalizeExprArray(exprArray *[]string) {
	if exprArray == nil {
		return
	}
	for idx, expr := range *exprArray {
		normalizeExpr(&expr)
		(*exprArray)[idx] = expr
	}
}

func normalizeGroupingSet(groupingSet planNodeGroupingSet) {
	if groupingSet.GroupKeys != nil {
		for idx, groupKey := range *groupingSet.GroupKeys {
			normalizeExprArray(&groupKey) // value comes from deparse_expression
			(*groupingSet.GroupKeys)[idx] = groupKey
		}
	}
	if groupingSet.HashKeys != nil {
		for idx, hashKey := range *groupingSet.HashKeys {
			normalizeExprArray(&hashKey) // value comes from deparse_expression
			(*groupingSet.HashKeys)[idx] = hashKey
		}
	}
	normalizeExprArray(groupingSet.SortKey) // value comes from deparse_expression
}

func normalizePlanNodeFields(node planNode) {
	normalizeExpr(node.CacheKey)       // value comes from deparse_expression
	normalizeExpr(node.ConflictFilter) // value comes from show_upper_qual
	normalizeExpr(node.Filter)         // value comes from show_scan_qual
	normalizeExpr(node.FunctionCall)   // value comes from show_expression
	normalizeExprArray(node.GroupKey)  // value comes from deparse_expression
	if node.GroupingSets != nil {
		for _, groupingSet := range *node.GroupingSets {
			normalizeGroupingSet(groupingSet)
		}
	}
	normalizeExpr(node.HashCond)                // value comes from show_upper_qual
	normalizeExpr(node.IndexCond)               // value comes from show_scan_qual
	normalizeExpr(node.JoinFilter)              // value comes from show_upper_qual
	normalizeExpr(node.MergeCond)               // value comes from show_upper_qual
	normalizeExpr(node.OneTimeFilter)           // value comes from show_upper_qual
	normalizeExpr(node.OrderBy)                 // value comes from show_scan_qual
	normalizeExprArray(node.Output)             // value comes from deparse_expression
	normalizeExprArray(node.PresortedKey)       // value comes from deparse_expression
	normalizeExpr(node.RecheckCond)             // value comes from show_scan_qual
	normalizeExpr(node.RepeatableSeed)          // value comes from deparse_expression
	normalizeExprArray(node.SamplingParameters) // value comes from deparse_expression
	if node.SortKey != nil {
		for idx, expr := range *node.SortKey {
			normalizeSortKey(&expr) // value comes from deparse_expression (but may have sort order information)
			(*node.SortKey)[idx] = expr
		}
	}
	normalizeExpr(node.TableFunctionCall) // value comes from show_expression
	normalizeExpr(node.TIDCond)           // value comes from show_scan_qual

	if node.Workers != nil {
		for _, p := range *node.Workers {
			normalizePlanNodeFields(p)
		}
	}

	for _, p := range node.Plans {
		normalizePlanNodeFields(p)
	}
}

func normalizePlan(planRaw json.RawMessage) (json.RawMessage, error) {
	var plan planNode
	var err error
	if err = json.Unmarshal(planRaw, &plan); err != nil {
		return planRaw, err
	}
	normalizePlanNodeFields(plan)
	return json.Marshal(plan)
}

// NormalizeExplainJSON - Normalizes the expressions contained within the
// passed in EXPLAIN JSON output.
func NormalizeExplainJSON(explainOutputJSON *state.ExplainPlanContainer) (*state.ExplainPlanContainer, error) {
	var err error
	explainOutputJSON.Plan, err = normalizePlan(explainOutputJSON.Plan)
	if err != nil {
		return nil, err
	}
	return explainOutputJSON, nil
}
