package ai

import (
	"sort"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00055 = "SQLE00055"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00055,
			Desc:       plocale.Rule00055Desc,
			Annotation: plocale.Rule00055Annotation,
			Category:   plocale.RuleTypeIndexOptimization,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00055Message,
		Func:    RuleSQLE00055,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00055): "For table creation and index creation statements, creating redundant indexes is prohibited .".
You should follow the following logic:
1. For the "CREATE TABLE ..." statements, builds a list of index columns, which is used to record all the declared indexes and their columns, checking whether the columns of each index are redundant, meaning that the index columns are exactly the same, or have the same leftmost prefix.  If it does, report a violation.
2. For the  "CREATE INDEX ..." statements, builds a list of index columns to keep track of the existing indexes and their columns, and check that the new index is not redundant, meaning that the index field is the same as the old one or has the same leftmost prefix. If it does, report a violation.
3. For the  "ALTER TABLE ... ADD INDEX ..." statements, perform the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00055(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."
		colIndexes := map[string][]string{}
		// 获取列定义中的索引列
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) {
				colIndexes["UNIQUE["+col.Name.Name.O+"]"] = []string{util.GetColumnName(col)}
			}
			if util.IsColumnPrimaryKey(col) {
				colIndexes["Primary Key"] = []string{util.GetColumnName(col)}
			}
		}

		// 获取表约束中的索引列
		constraintsIndexes := extractIndexesFromConstraints(util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...))
		// 此处要聚合到一个map中，否则无法获取存在于表约束中的冗余索引
		mergeMap := mergeIndexsMaps(colIndexes, constraintsIndexes)
		// 获取冗余索引
		existsIndexs, redundantIndexs := GroupDuplicatesByValue(mergeMap)
		buildAuditResult(input, existsIndexs, redundantIndexs)

	case *ast.CreateIndexStmt:
		// "create index..."

		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		indexes := extractIndexesFromConstraints(util.GetTableConstraints(createTableStmt.Constraints, util.GetIndexConstraintTypes()...))
		newIndexMap := make(map[string][]string)
		newIndexMap[stmt.IndexName] = extractIndexesFromIndexStmt(stmt.IndexPartSpecifications)
		// 获取冗余索引
		existsIndexs, redundantIndexs := GetIntersectionByValue(indexes, newIndexMap)
		buildAuditResult(input, existsIndexs, redundantIndexs)

	case *ast.AlterTableStmt:
		// "alter table"

		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		indexes := extractIndexesFromConstraints(util.GetTableConstraints(createTableStmt.Constraints, util.GetIndexConstraintTypes()...))
		constraints := make([]*ast.Constraint, 0)
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table... add index..."
			constraints = append(constraints, spec.Constraint)
		}

		// // 获取冗余索引
		newIndexs := extractIndexesFromConstraints(util.GetTableConstraints(constraints, util.GetIndexConstraintTypes()...))
		existsIndexs, redundantIndexs := GetIntersectionByValue(indexes, newIndexs)
		buildAuditResult(input, existsIndexs, redundantIndexs)
	}

	return nil
}

func buildAuditResult(input *rulepkg.RuleHandlerInput, existsIndexs, redundantIndexs map[string][]string) {
	idxNames := make([]string, 0, len(existsIndexs))
	idxColNames := make([][]string, 0, len(existsIndexs))
	redundantIdxNames := make([]string, 0, len(redundantIndexs))
	for key, idx := range existsIndexs {
		idxNames = append(idxNames, key)
		idxColNames = append(idxColNames, idx)
	}
	for key := range redundantIndexs {
		redundantIdxNames = append(redundantIdxNames, key)
	}
	if len(idxNames) > 0 && len(idxColNames) > 0 && len(redundantIdxNames) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00055, idxColNames, strings.Join(idxNames, "、"), redundantIdxNames)
	}
}

func extractIndexesFromConstraints(constraints []*ast.Constraint) map[string][]string {
	indexes := make(map[string][]string)

	// Iterate over all constraints to extract index columns.
	for _, constraint := range constraints {
		indexCols := []string{}

		// Collect all column names that are part of the index represented by the constraint.
		for _, key := range constraint.Keys {
			colName := util.GetIndexColName(key)
			if colName != "" {
				indexCols = append(indexCols, colName)
			}
		}

		// Only append to indexes if indexCols is not empty, avoiding adding empty index definitions.
		if len(indexCols) > 0 {
			key := constraint.Name
			if key == "" && constraint.Tp == ast.ConstraintPrimaryKey {
				key = "Primary Key"
			}
			indexes[key] = indexCols
		}
	}

	return indexes
}

func mergeIndexsMaps(m1, m2 map[string][]string,
) map[string][]string {
	result := make(map[string][]string, len(m1)+len(m2))
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result
}

func extractIndexesFromIndexStmt(index []*ast.IndexPartSpecification) []string {
	// Initialize the slice that will hold the column names for the index.
	var indexCols []string

	// Iterate over each IndexPartSpecification to extract the column names.
	for _, col := range index {
		if colName := util.GetIndexColName(col); colName != "" {
			indexCols = append(indexCols, colName)
		}
	}

	// Return the indexes as a slice of a single index if indexCols is not empty.
	if len(indexCols) > 0 {
		return indexCols
	}
	return []string{} // Return an empty slice of indexes if no index columns were found.
}

// GetIntersectionByValue 比较两个映射的值集合，返回两个结果映射：
//
//	commonInLeft:  左映射中，其值集合（或其左前缀）出现在右映射中的条目
//	commonInRight: 右映射中，其值集合（或其左前缀）出现在左映射中的条目
//
// 参数：
//
//	leftMap  - 键到字符串切片的映射，作为左侧数据源
//	rightMap - 键到字符串切片的映射，作为右侧数据源
//
// 返回：
//
//	commonInLeft  - 所有 leftMap 中，其值（或左前缀）能在 rightMap 中找到的项
//	commonInRight - 所有 rightMap 中，其值（或左前缀）能在 leftMap 中找到的项
func GetIntersectionByValue(leftMap, rightMap map[string][]string) (commonInLeft map[string][]string, commonInRight map[string][]string) {

	// 初始化返回容器
	commonInLeft = make(map[string][]string)
	commonInRight = make(map[string][]string)

	// --- 构建 rightFullIndex：右侧完整切片的快速查找索引 ---
	// 示例：rightMap 中有 entry ["a","b","c"]，则索引中添加 "a,b,c"
	rightFullIndex := make(map[string]bool)
	for _, values := range rightMap {
		fullKey := strings.Join(values, ",")
		rightFullIndex[fullKey] = true
	}

	// --- 扫描 leftMap，检测左侧切片或其左前缀是否命中 rightFullIndex ---
	for key, values := range leftMap {
		// 从最长前缀到最短前缀依次尝试
		for length := len(values); length > 0; length-- {
			prefixKey := strings.Join(values[:length], ",")
			if rightFullIndex[prefixKey] {
				// 一旦命中，记录该 leftMap 条目为 commonInLeft，跳出前缀循环
				commonInLeft[key] = copySlice(values)
				break
			}
		}
	}

	// --- 构造 leftPrefixIndex：左侧所有可能前缀的索引 ---
	// 为了对称地支持右侧完整切片的匹配，需要将左侧每个 values 的所有左前缀都注册
	// 例如 values=["a","b","c"]，则注册 "a"，"a,b"，"a,b,c"
	leftPrefixIndex := make(map[string]bool)
	for _, values := range leftMap {
		for i := 1; i <= len(values); i++ {
			sig := strings.Join(values[:i], ",")
			leftPrefixIndex[sig] = true
		}
	}

	// --- 扫描 rightMap，检测右侧完整切片是否命中 leftPrefixIndex ---
	for key, values := range rightMap {
		fullKey := strings.Join(values, ",")
		if leftPrefixIndex[fullKey] {
			// 如果右侧完整切片恰好是左侧某个切片的左前缀，则记录为 commonInRight
			commonInRight[key] = copySlice(values)
		}
	}

	return
}

// GroupDuplicatesByValue 检测值集合重复的键，返回分组结果：
// - originals: 每个唯一值集合（或其最长左前缀）首次出现的键及其值
// - duplicates: 后续出现且与某已登记左前缀匹配的键及其值
func GroupDuplicatesByValue(inputMap map[string][]string) (originals map[string][]string, duplicates map[string][]string) {
	originals = make(map[string][]string)
	duplicates = make(map[string][]string)

	// prefixIndex：将“值切片左前缀签名”映射到首次出现该前缀的 key
	prefixIndex := make(map[string]string)

	// 为保证稳定性，先对所有 key 排序
	keys := make([]string, 0, len(inputMap))
	for k := range inputMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := inputMap[key]
		matched := false

		// 从最长前缀到最短前缀尝试匹配
		for length := len(values); length > 0; length-- {
			sig := strings.Join(values[:length], ",")
			if origKey, exists := prefixIndex[sig]; exists {
				// 找到已有前缀，归为 duplicates
				if !matched {
					// 确保 originals 中保留最早出现的完整值
					if _, ok := originals[origKey]; !ok {
						originals[origKey] = copySlice(inputMap[origKey])
					}
					matched = true
				}
				duplicates[key] = copySlice(values)
				break
			}
		}
		if !matched {
			// 全新起点，将自身所有左前缀注册
			for length := len(values); length > 0; length-- {
				sig := strings.Join(values[:length], ",")
				prefixIndex[sig] = key
			}
		}
	}

	return
}

// copySlice 深拷贝字符串slice
func copySlice(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// ==== Rule code end ====
