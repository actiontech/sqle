//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/knowledge_base"
	"github.com/labstack/echo/v4"
)

// 获取知识库列表
func getKnowledgeBaseList(c echo.Context) error {
	var req GetKnowledgeBaseListReq
	if err := controller.BindAndValidateReq(c, &req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	knowledgeList, count, err := knowledge_base.SearchKnowledgeList(req.KeyWords, req.Tags, int(limit), int(offset))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, GetKnowledgeBaseListRes{
		BaseRes:   controller.NewBaseReq(nil),
		TotalNums: uint64(count),
		Data:      convertToKnowledgeBaseListRes(knowledgeList),
	})
}

func convertToKnowledgeBaseListRes(knowledgeList []model.Knowledge) []*KnowledgeBase {
	knowledgeRes := make([]*KnowledgeBase, 0, len(knowledgeList))
	for _, knowledge := range knowledgeList {
		knowledgeRes = append(knowledgeRes, &KnowledgeBase{
			ID:          knowledge.ID,
			Title:       knowledge.Title,
			Description: knowledge.Description,
			Content:     knowledge.Content,
			Tags:        convertToTagRes(knowledge.Tags),
		})
	}
	return knowledgeRes
}

func convertToTagRes(tags []*model.Tag) []*Tag {
	tagRes := make([]*Tag, 0, len(tags))
	for _, tag := range tags {
		subTags := make([]*Tag, 0, len(tag.SubTag))
		for _, subTag := range tag.SubTag {
			subTags = append(subTags, &Tag{
				ID:   subTag.ID,
				Name: fmt.Sprint(subTag.Name),
			})
		}
		tagRes = append(tagRes, &Tag{
			ID:      tag.ID,
			Name:    fmt.Sprint(tag.Name),
			SubTags: subTags,
		})
	}
	return tagRes
}

// 获取知识库标签列表
func getKnowledgeBaseTagList(c echo.Context) error {
	tags, err := knowledge_base.GetKnowledgeBaseTags()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	tagRes := convertToTagRes(tags)

	return c.JSON(http.StatusOK, GetKnowledgeBaseTagListRes{
		BaseRes:   controller.NewBaseReq(nil),
		TotalNums: uint64(len(tagRes)),
		Data:      tagRes,
	})
}
