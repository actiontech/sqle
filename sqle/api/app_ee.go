//go:build enterprise
// +build enterprise

package api

import (
	"net/http"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
)

func init() {
	LoadRestApi(http.MethodGet, "/projects/names", v1.GetProjectNamesByIds) // 浙农信定制
}
