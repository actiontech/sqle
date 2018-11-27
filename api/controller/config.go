package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"sqle/inspector"
	"sqle/model"
)

type GetAllConfigRes struct {
	BaseRes
	Data []model.Config `json:"data"`
}

// @Summary 配置列表
// @Description get all configs
// @Success 200 {object} controller.GetAllConfigRes
// @router /configs [get]
func GetAllConfig(c echo.Context) error {
	s := model.GetStorage()
	configs, err := s.GetAllConfig()
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	return c.JSON(200, &GetAllConfigRes{
		BaseRes: NewBaseReq(nil),
		Data:    configs,
	})
}

type UpdateConfigReq struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UpdateAllConfigReq struct {
	Configs []UpdateConfigReq `json:"config_list"`
}

// @Summary 修改配置
// @Description update configs
// @Accept json
// @Produce json
// @Param instance body controller.UpdateAllConfigReq true "update config"
// @Success 200 {object} controller.BaseRes
// @router /configs [patch]
func UpdateConfigs(c echo.Context) error {
	s := model.GetStorage()
	reqs := new(UpdateAllConfigReq)
	if err := c.Bind(reqs); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	configs, err := s.GetConfigMap()
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	for _, req := range reqs.Configs {
		if _, ok := configs[req.Name]; ok {
			err := s.UpdateConfigValueByName(req.Name, req.Value)
			if err != nil {
				return c.JSON(http.StatusOK, NewBaseReq(err))
			}
			inspector.UpdateConfig(req.Name, req.Value)
		}
	}
	return c.JSON(200, NewBaseReq(nil))
}
