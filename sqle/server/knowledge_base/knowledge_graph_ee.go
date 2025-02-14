//go:build enterprise
// +build enterprise

package knowledge_base

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

// 节点
type Node struct {
	ID     string      // 节点ID
	Name   string      // 节点名称
	Weight uint64      // 权重
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
			return
		}
	}

	// 如果没有找到相同的边，添加新边
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
func GetGraph() (*Graph, error) {
	builder := NewGraphBuilder(model.GetStorage())
	return builder.Build()
}

// Build constructs the complete knowledge graph.
func (b *GraphBuilder) Build() (*Graph, error) {
	// Fetch all required data
	tagData, err := b.fetchTagData()
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
	fatherTag             *model.Tag
	tags                  []*model.Tag
	knowledgeTagRelations []*model.KnowledgeTagRelation
}

// graphData holds the processed data ready for graph construction
type graphData struct {
	tagMap          map[uint]*model.Tag
	tagCountMap     map[uint]uint64
	knowledgeTagMap map[uint][]uint
	edgeWeightMap   map[string]uint64 // 存储边的权重
}

// edgeKey generates a unique key for an edge between two tag IDs
func makeEdgeKey(fromID, toID uint) string {
	if fromID < toID {
		return fmt.Sprintf("%d:%d", fromID, toID)
	}
	return fmt.Sprintf("%d:%d", toID, fromID)
}

// fetchTagData retrieves all necessary data from storage
func (b *GraphBuilder) fetchTagData() (*tagData, error) {
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
	relations, err := b.storage.GetKnowledgeBaseTagRelations()
	if err != nil {
		log.Logger().Errorf("failed to get knowledge base tag relations: %v", err)
		return nil, err
	}

	return &tagData{
		fatherTag:             fatherTag,
		tags:                  tags,
		knowledgeTagRelations: relations,
	}, nil
}

// transformData processes raw data into formats ready for graph construction
func (b *GraphBuilder) transformData(data *tagData) (*graphData, error) {
	// Pre-allocate maps with estimated capacity
	tagCount := len(data.tags)
	relationCount := len(data.knowledgeTagRelations)

	tagMap := make(map[uint]*model.Tag, tagCount)
	tagCountMap := make(map[uint]uint64, tagCount)
	knowledgeTagMap := make(map[uint][]uint, relationCount)
	edgeWeightMap := make(map[string]uint64, relationCount)

	// Build tag map for quick lookups
	for i := range data.tags {
		tagMap[data.tags[i].ID] = data.tags[i]
		// Initialize tag count to 1 as default weight
		tagCountMap[data.tags[i].ID] = 1
	}

	// Count tag occurrences and build knowledge-tag relationships
	for _, relation := range data.knowledgeTagRelations {
		tagCountMap[relation.TagID]++
		knowledgeTagMap[relation.KnowledgeID] = append(
			knowledgeTagMap[relation.KnowledgeID],
			relation.TagID,
		)
	}

	// Calculate edge weights
	for _, tagIDs := range knowledgeTagMap {
		for i := 0; i < len(tagIDs)-1; i++ {
			edgeKey := makeEdgeKey(tagIDs[i], tagIDs[i+1])
			// Initialize edge weight to 1 if not exists, otherwise increment
			edgeWeightMap[edgeKey]++
		}
	}

	return &graphData{
		tagMap:          tagMap,
		tagCountMap:     tagCountMap,
		knowledgeTagMap: knowledgeTagMap,
		edgeWeightMap:   edgeWeightMap,
	}, nil
}

// constructGraph builds the final graph structure from processed data
func (b *GraphBuilder) constructGraph(data *graphData) (*Graph, error) {
	graph := &Graph{
		Nodes: make([]*Node, 0, len(data.tagMap)),
		Edges: make([]*Edge, 0, len(data.knowledgeTagMap)*2), // Estimate edge count
	}

	// Create node map for edge construction
	nodeMap := make(map[uint]*Node, len(data.tagMap))

	// Add all nodes first
	for tagID, tag := range data.tagMap {
		node := &Node{
			ID:     fmt.Sprint(tagID),
			Name:   string(tag.Name),
			Weight: data.tagCountMap[tagID], // Weight starts from 1 and increases with occurrences
		}
		graph.AddNode(node)
		nodeMap[tagID] = node
	}

	// Add edges based on knowledge relationships
	for _, tagIDs := range data.knowledgeTagMap {
		for i := 0; i < len(tagIDs)-1; i++ {
			edgeKey := makeEdgeKey(tagIDs[i], tagIDs[i+1])
			edge := &Edge{
				From:       nodeMap[tagIDs[i]],
				To:         nodeMap[tagIDs[i+1]],
				Weight:     data.edgeWeightMap[edgeKey], // Set weight based on occurrence count
				IsDirected: true,
			}
			graph.AddEdge(edge)
		}
	}

	return graph, nil
}
