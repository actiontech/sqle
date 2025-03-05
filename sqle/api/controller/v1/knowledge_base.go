package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetKnowledgeBaseListReq struct {
	KeyWords  string   `json:"keywords" query:"keywords" example:"keywords"`   // 搜索内容
	Tags      []string `json:"tags" query:"tags" example:"tag1"`               // 搜索标签
	PageIndex uint32   `json:"page_index" query:"page_index" valid:"required"` // 页码
	PageSize  uint32   `json:"page_size" query:"page_size" valid:"required"`   // 每页条数
}

type GetKnowledgeBaseListRes struct {
	controller.BaseRes
	Data      []*KnowledgeBase `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

type KnowledgeBase struct {
	ID          uint   `json:"id"`          // 知识库ID
	RuleName    string `json:"rule_name"`   // 规则名称
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述
	Content     string `json:"content"`     // 内容
	Tags        []*Tag `json:"tags"`        // 标签
}

type Tag struct {
	ID      uint   `json:"id"`                 // 标签ID
	Name    string `json:"name"`               // 标签名称
	SubTags []*Tag `json:"sub_tags,omitempty"` // 子标签
}

// GetKnowledgeBaseList
// @Summary 获取知识库列表
// @Description get knowledge base list
// @Id getKnowledgeBaseList
// @Tags knowledge_base
// @Param keywords query string false "keywords"
// @Param tags query []string false "tags"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetKnowledgeBaseListRes
// @router /v1/knowledge_bases [get]
func GetKnowledgeBaseList(c echo.Context) error {
	return getKnowledgeBaseList(c)
}

type GetKnowledgeBaseTagListRes struct {
	controller.BaseRes
	TotalNums uint64 `json:"total_nums"`
	Data      []*Tag `json:"data"`
}

// GetKnowledgeBaseTagList
// @Summary 获取知识库标签列表
// @Description get tag list of knowledge base
// @Id getKnowledgeBaseTagList
// @Tags knowledge_base
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetKnowledgeBaseTagListRes
// @router /v1/knowledge_bases/tags [get]
func GetKnowledgeBaseTagList(c echo.Context) error {
	return getKnowledgeBaseTagList(c)
}

// GetKnowledgeGraph
// @Summary 获取知识库知识图谱
// @Description get knowledge graph
// @Id getKnowledgeGraph
// @Tags knowledge_base
// @Param filter_by_rule_name query string false "filter by rule name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetKnowledgeGraphResp
// @router /v1/knowledge_bases/graph [get]
func GetKnowledgeGraph(c echo.Context) error {
	return getKnowledgeGraph(c)
}

// 获取知识图谱请求
type GetKnowledgeGraphReq struct {
	FilterByRuleName string `json:"filter_by_rule_name" query:"filter_by_rule_name"` // 根据规则名称过滤
}

type GetKnowledgeGraphResp struct {
	controller.BaseRes
	Data *GraphResponse `json:"data"`
}

// NodeResponse represents a node in the API response
type NodeResponse struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Weight uint64      `json:"weight"`
}

// EdgeResponse represents an edge in the API response
type EdgeResponse struct {
	FromID     string `json:"from_id"`     // 存储Node的ID
	ToID       string `json:"to_id"`       // 存储Node的ID
	Weight     uint64 `json:"weight"`      // 权重
	IsDirected bool   `json:"is_directed"` // 是否有向
}

// GraphResponse represents the complete graph structure in the API response
type GraphResponse struct {
	Nodes []*NodeResponse `json:"nodes"` // 节点集合
	Edges []*EdgeResponse `json:"edges"` // 边集合
}
