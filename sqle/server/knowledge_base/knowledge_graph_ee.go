package knowledge_base

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

// 节点
type Node struct {
	ID     string // 节点ID
	Name   string // 节点名称
	Weight uint64 // 权重
}

// 边
type Edge struct {
	From       *Node  // 起点
	To         *Node  // 终点
	Weight     uint64 // 权重
	IsDirected bool   // 是否有向
}

// 图
type Graph struct {
	Nodes []*Node // 节点集合
	Edges []*Edge // 边集合
}

// 判断节点是否已存在
func (g *Graph) hasNode(node *Node) bool {
	if node == nil {
		return false
	}
	for _, existingNode := range g.Nodes {
		if existingNode != nil && existingNode.ID == node.ID {
			return true
		}
	}
	return false
}

// 添加节点到图，如果节点已存在则不添加
func (g *Graph) AddNode(node *Node) {
	if node == nil {
		return
	}
	if !g.hasNode(node) {
		g.Nodes = append(g.Nodes, node)
	}
}

// 判断两条边是否相同
func (g *Graph) isSameEdge(e1, e2 *Edge) bool {
	if e1.IsDirected {
		// 有向图：起点和终点都要相同
		return e1.From == e2.From && e1.To == e2.To
	}
	// 无向图：起点和终点可以互换
	return (e1.From == e2.From && e1.To == e2.To ||
		e1.From == e2.To && e1.To == e2.From)
}

// 添加边到图，如果边已存在则增加权重
func (g *Graph) AddEdge(edge *Edge) {
	if edge == nil {
		return
	}
	if edge.From == nil || edge.To == nil {
		return
	}
	// 如果起点和终点不在图中则不添加
	if !g.hasNode(edge.From) || !g.hasNode(edge.To) {
		return
	}

	// 查找是否存在相同的边
	for _, existingEdge := range g.Edges {
		if g.isSameEdge(existingEdge, edge) {
			// 如果存在相同的边，增加权重
			existingEdge.Weight++
			return
		}
	}

	// 如果没有找到相同的边，添加新边
	// 初始权重为1
	edge.Weight = 1
	g.Edges = append(g.Edges, edge)
}

// GraphBuilder handles the construction of knowledge graphs from tags and their relationships.
type GraphBuilder struct {
	storage *model.Storage
}

// NewGraphBuilder creates a new GraphBuilder instance.
func NewGraphBuilder(storage *model.Storage) *GraphBuilder {
	return &GraphBuilder{
		storage: storage,
	}
}

// GetGraph constructs and returns a knowledge graph based on tags and their relationships.
// The graph represents the connections between different knowledge tags, where:
// - Nodes represent individual tags with default weight of 1
// - Edges represent relationships between tags based on knowledge base connections
// - Node weights increase based on tag occurrence frequency
// - Edge weights increase when the same connection appears multiple times
func GetGraph(ruleName string) (*Graph, error) {
	builder := NewGraphBuilder(model.GetStorage())
	return builder.Build(ruleName)
}

// Build constructs the complete knowledge graph.
func (b *GraphBuilder) Build(ruleName string) (*Graph, error) {
	// Fetch all required data
	tagData, err := b.fetchTagData(ruleName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tag data: %w", err)
	}

	// Transform data into intermediate format
	graphData, err := b.transformData(tagData)
	if err != nil {
		return nil, fmt.Errorf("failed to transform data: %w", err)
	}

	// Construct the final graph
	return b.constructGraph(graphData)
}

// tagData holds all the raw data needed to build the graph
type tagData struct {
	fatherTag        *model.Tag
	tags             []*model.Tag
	knowledgeTagSets map[string][][]model.TypeTag
}

// graphData holds the processed data ready for graph construction
type graphData struct {
	tagMap           map[model.TypeTag]*model.Tag
	tagCountMap      map[uint]uint64
	knowledgeTagSets map[string][][]model.TypeTag
}

// fetchTagData retrieves all necessary data from storage
func (b *GraphBuilder) fetchTagData(ruleName string) (*tagData, error) {
	// Get knowledge base predefined tag
	fatherTag, err := b.storage.GetTagByName(model.PredefineTagKnowledgeBase)
	if err != nil {
		log.Logger().Errorf("failed to get father tag: %v", err)
		return nil, err
	}

	// Get all sub tags
	tags, err := b.storage.GetSubTags(fatherTag.ID)
	if err != nil {
		log.Logger().Errorf("failed to get sub tags: %v", err)
		return nil, err
	}

	// Get knowledge-tag relations
	var knowledgeTagSets map[string][][]model.TypeTag
	knowledgeTagSets = model.GetTagMapDefaultRuleKnowledge()
	if err != nil {
		log.Logger().Errorf("failed to get knowledge base tag relations: %v", err)
		return nil, err
	}
	if ruleName != "" {
		knowledgeTagSets = map[string][][]model.TypeTag{ruleName: knowledgeTagSets[ruleName]}
	}

	return &tagData{
		fatherTag:        fatherTag,
		tags:             tags,
		knowledgeTagSets: knowledgeTagSets,
	}, nil
}

// transformData processes raw data into formats ready for graph construction
func (b *GraphBuilder) transformData(data *tagData) (*graphData, error) {
	// Pre-allocate maps with estimated capacity
	tagCount := len(data.tags)

	tagMap := make(map[model.TypeTag]*model.Tag, tagCount)
	tagCountMap := make(map[uint]uint64, tagCount)

	// Build tag map for quick lookups
	for i := range data.tags {
		tagMap[data.tags[i].Name] = data.tags[i]
		// Initialize tag count to 1 as default weight
		tagCountMap[data.tags[i].ID] = 1
	}

	for _, tagSets := range model.GetTagMapDefaultRuleKnowledge() {
		for _, tagSet := range tagSets {
			for _, tagName := range tagSet {
				tagCountMap[tagMap[tagName].ID]++
			}
		}
	}

	return &graphData{
		tagMap:           tagMap,
		tagCountMap:      tagCountMap,
		knowledgeTagSets: data.knowledgeTagSets,
	}, nil
}

func GetObjectNodesMap() map[model.TypeTag]*model.Tag {
	return map[model.TypeTag]*model.Tag{
		model.TABLE:       nil,
		model.INDEX:       nil,
		model.COLUMN:      nil,
		model.VIEW:        nil,
		model.DATABASE:    nil,
		model.SCHEMA:      nil,
		model.PROCEDURE:   nil,
		model.FUNCTION:    nil,
		model.TRIGGER:     nil,
		model.USER:        nil,
		model.ROLE:        nil,
		model.EVENT:       nil,
		model.PARTITION:   nil,
		model.SEQUENCE:    nil,
		model.KEY:         nil,
		model.ROW:         nil,
		model.VALUES:      nil,
		model.STATEMENT:   nil,
		model.TRANSACTION: nil,
		model.SUBQUERY:    nil,
	}
}

func GetObjectDescribeNodesMap() map[model.TypeTag]*model.Tag {
	return map[model.TypeTag]*model.Tag{
		model.UNIQUE:         nil,
		model.PRIMARY:        nil,
		model.FOREIGN:        nil,
		model.ENUM:           nil,
		model.VARCHAR:        nil,
		model.CHAR:           nil,
		model.DECIMAL:        nil,
		model.BLOB:           nil,
		model.INT:            nil,
		model.NOT_NULL:       nil,
		model.NULL:           nil,
		model.EXISTS:         nil,
		model.COMMENT:        nil,
		model.CHARSET:        nil,
		model.COLLATION:      nil,
		model.ENGINE:         nil,
		model.MODIFY:         nil,
		model.LENGTH:         nil,
		model.CONSTRAINT:     nil,
		model.TIMESTAMP:      nil,
		model.RAND:           nil,
		model.FLOAT:          nil,
		model.BIGINT:         nil,
		model.IS:             nil,
		model.AUTO_INCREMENT: nil,
		model.NOT:            nil,
		model.LEVEL:          nil,
		model.GLOBAL:         nil,
	}
}

func GetManipulateNodesMap() map[model.TypeTag]*model.Tag {
	return map[model.TypeTag]*model.Tag{
		model.SELECT:     nil,
		model.INSERT:     nil,
		model.UPDATE:     nil,
		model.DELETE:     nil,
		model.CREATE:     nil,
		model.ALTER:      nil,
		model.DROP:       nil,
		model.RENAME:     nil,
		model.ADD:        nil,
		model.SET:        nil,
		model.LIMIT:      nil,
		model.IN:         nil,
		model.AS:         nil,
		model.WHERE:      nil,
		model.FROM:       nil,
		model.CONSTRAINT: nil,
		model.ON:         nil,
		model.ORDER:      nil,
		model.GROUP:      nil,
		model.HAVING:     nil,
		model.OR:         nil,
		model.SUBQUERY:   nil,
		model.UNION:      nil,
		model.JOIN:       nil,
		model.OFFSET:     nil,
		model.BY:         nil,
		model.TO:         nil,
		model.IF:         nil,
		model.HINT:       nil,
		model.EXPLAIN:    nil,
		model.COUNT:      nil,
		model.ISOLATION:  nil,
		model.PRIVILEGE:  nil,
		model.TEMPORARY:  nil,
		model.FOR:        nil,
		model.LIKE:       nil,
		model.EACH:       nil,
		model.INTO:       nil,
		model.TRUNCATE:   nil,
		model.GRANT:      nil,
		model.REFERENCES: nil,
	}
}

func (b *GraphBuilder) constructGraph(data *graphData) (*Graph, error) {
	graph := &Graph{
		Nodes: make([]*Node, 0, len(data.tagMap)),
		Edges: make([]*Edge, 0),
	}

	// Create node map for edge construction
	nodeMap := make(map[uint]*Node, len(data.tagMap))

	// Add all nodes first
	for _, tag := range data.tagMap {
		node := &Node{
			ID:     fmt.Sprint(tag.ID),
			Name:   string(tag.Name),
			Weight: data.tagCountMap[tag.ID], // Weight starts from 1 and increases with occurrences
		}
		nodeMap[tag.ID] = node
	}
	// Add nodes based on knowledge-tag relationships
	for _, tagSets := range data.knowledgeTagSets {
		for _, tagSet := range tagSets {
			for _, tagName := range tagSet {
				if tag, exist := data.tagMap[tagName]; exist {
					graph.AddNode(nodeMap[tag.ID])
				}
			}
		}
	}

	// Add edges based on knowledge-tag relationships
	for _, tagSets := range data.knowledgeTagSets {
		var objectNodeIDMap = make(map[uint]struct{})
		for _, tagSet := range tagSets {
			// 找出中心节点、边缘节点和中间节点
			var centerNodeIDs []uint
			var edgeNodeIDs []uint
			var middleNodeIDs []uint
			for _, tagName := range tagSet {
				if _, exist := GetObjectNodesMap()[tagName]; exist {
					centerNodeIDs = append(centerNodeIDs, data.tagMap[tagName].ID)
					objectNodeIDMap[data.tagMap[tagName].ID] = struct{}{}
				}
				if _, exist := GetManipulateNodesMap()[tagName]; exist {
					edgeNodeIDs = append(edgeNodeIDs, data.tagMap[tagName].ID)
				}
				if _, exist := GetObjectDescribeNodesMap()[tagName]; exist {
					middleNodeIDs = append(middleNodeIDs, data.tagMap[tagName].ID)
				}
			}
			// 创建由中心节点到中间节点的边
			for _, centerNodeID := range centerNodeIDs {
				for _, middleNodeID := range middleNodeIDs {
					edge := &Edge{
						From:       nodeMap[centerNodeID],
						To:         nodeMap[middleNodeID],
						IsDirected: false,
					}
					graph.AddEdge(edge)
				}
			}
			// 创建由中心节点到边缘节点的边
			for _, centerNodeID := range centerNodeIDs {
				for _, edgeNodeID := range edgeNodeIDs {
					edge := &Edge{
						From:       nodeMap[centerNodeID],
						To:         nodeMap[edgeNodeID],
						IsDirected: false,
					}
					graph.AddEdge(edge)
				}
			}
		}
		// 如果有多个Object则连接每一个Object
		if len(objectNodeIDMap) > 1 {
			objectNodeIDs := make([]uint, 0, len(objectNodeIDMap))
			for objectNodeID := range objectNodeIDMap {
				objectNodeIDs = append(objectNodeIDs, objectNodeID)
			}
			// 每2个object构成一条边
			for i := 0; i < len(objectNodeIDs); i++ {
				for j := i + 1; j < len(objectNodeIDs); j++ {
					edge := &Edge{
						From:       nodeMap[objectNodeIDs[i]],
						To:         nodeMap[objectNodeIDs[j]],
						IsDirected: false,
					}
					graph.AddEdge(edge)
				}
			}
		}
	}
	return graph, nil
}
