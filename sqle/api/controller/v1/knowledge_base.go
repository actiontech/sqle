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
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述
	Content     string `json:"content"`     // 内容
	Tags        []*Tag `json:"tags"`        // 标签
}

type Tag struct {
	ID   uint   `json:"id"`   // 标签ID
	Name string `json:"name"` // 标签名称
}

// GetKnowledgeBaseList
// @Summary 获取知识库列表
// @Description get knowledge base list
// @Id getKnowledgeBaseList
// @Tags knowledge_base
// @Param keywords query string false "keywords"
// @Param tags query []string false "tags"
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
