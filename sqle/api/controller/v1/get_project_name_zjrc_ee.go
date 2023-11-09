//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type GetProjectNameQuery struct {
	Ids []int64 `query:"ids"` // 浙农信系统中的项目id
}

type GetProjectNameResV1 struct {
	controller.BaseRes
	Data []IDProjectName `json:"data"`
}

type IDProjectName struct {
	OriginProjectID int64  `json:"origin_project_id"` // 浙农信系统中的项目id
	SqleProjectName string `json:"sqle_project_name"` // sqle系统中的项目名称
}

func getProjectNamesByIdsResponse(pair []model.IdProjectNamePair) *GetProjectNameResV1 {
	data := make([]IDProjectName, 0, len(pair))
	for i := range pair {
		data = append(data, IDProjectName{
			OriginProjectID: pair[i].OriginProjectID,
			SqleProjectName: pair[i].SqleProjectName,
		})
	}
	return &GetProjectNameResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	}
}

func getProjectNamesByIds(c echo.Context) error {
	var query GetProjectNameQuery
	if err := controller.BindAndValidateReq(c, &query); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	idProjectNamePairs, _, err := model.GetStorage().GetIdProjectNamePairsByIds(query.Ids)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, getProjectNamesByIdsResponse(idProjectNamePairs))
}

func GetProjectNamesByIds(c echo.Context) error {
	return getProjectNamesByIds(c)
}
