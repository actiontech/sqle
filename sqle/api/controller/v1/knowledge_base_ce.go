//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionDoesNotSupportKnowledgeBase = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support knowledge base"))

func getKnowledgeBaseList(c echo.Context) error {
	return ErrCommunityEditionDoesNotSupportKnowledgeBase
}

func getKnowledgeBaseTagList(c echo.Context) error {
	return ErrCommunityEditionDoesNotSupportKnowledgeBase
}

func getKnowledgeGraph(c echo.Context) error {
	return ErrCommunityEditionDoesNotSupportKnowledgeBase
}