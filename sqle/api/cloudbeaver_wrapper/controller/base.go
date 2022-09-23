package controller

import (
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/resolver"

	"github.com/labstack/echo/v4"
)

type Next func(c echo.Context) ([]byte, error)

type ResolverImpl struct {
	*resolver.Resolver
	Ctx  echo.Context
	Next Next
}

func (r *ResolverImpl) Mutation() resolver.MutationResolver {
	return &MutationResolverImpl{
		Ctx:  r.Ctx,
		Next: r.Next,
	}
}

// Query returns generated.QueryResolver implementation.
func (r *ResolverImpl) Query() resolver.QueryResolver {
	return &QueryResolverImpl{
		Ctx:  r.Ctx,
		Next: r.Next,
	}
}

type MutationResolverImpl struct {
	*resolver.MutationResolverImpl
	Ctx  echo.Context
	Next Next
}

type QueryResolverImpl struct {
	*resolver.QueryResolverImpl
	Ctx  echo.Context
	Next Next
}
