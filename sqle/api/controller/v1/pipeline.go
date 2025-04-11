package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/server/pipeline"
	"github.com/labstack/echo/v4"
)

// pipelineDetail 流水线的信息详情
type pipelineDetail struct {
	ID        uint   `json:"id"`         // 流水线的唯一标识符
	NodeCount uint32 `json:"node_count"` // 节点个数
	pipelineBase
}

func (p *pipelineDetail) fillWith(pipe *pipeline.Pipeline) {
	p.ID = pipe.ID
	p.NodeCount = pipe.NodeCount()
	p.Name = pipe.Name
	p.Description = pipe.Description
	p.Address = pipe.Address
}

// pipelineBase 流水线基础信息
type pipelineBase struct {
	Name        string `json:"name"  valid:"required"` // 流水线名称
	Description string `json:"description"`            // 流水线描述
	Address     string `json:"address"`                // 关联流水线地址
}

func (pipe *pipelineBase) toSvcPipeline() *pipeline.Pipeline {
	return &pipeline.Pipeline{
		Name:        pipe.Name,
		Description: pipe.Description,
		Address:     pipe.Address,
	}
}

// pipelineNodeDetail 流水线节点的信息详情
type pipelineNodeDetail struct {
	ID              uint   `json:"id"`               // 节点的唯一标识符
	IntegrationInfo string `json:"integration_info"` // 对接说明
	pipelineNodeBase
}

func (p *pipelineNodeDetail) fillWith(ctx context.Context, node *pipeline.PipelineNode, projectName string) {
	if node == nil {
		return
	}
	integrationInfo, err := node.IntegrationInfo(ctx, projectName)
	if err != nil {
		integrationInfo = err.Error()
	}
	p.ID = node.ID
	p.Name = node.Name
	p.Type = node.NodeType
	p.InstanceName = node.InstanceName
	p.InstanceType = node.InstanceType
	p.ObjectPath = node.ObjectPath
	p.ObjectType = node.ObjectType
	p.AuditMethod = node.AuditMethod
	p.RuleTemplateName = node.RuleTemplateName
	p.IntegrationInfo = integrationInfo
}

// pipelineNodeBase 流水线节点基础信息
type pipelineNodeBase struct {
	Name             string `json:"name" valid:"required"`                                         // 节点名称，必填，支持中文、英文+数字+特殊字符
	Type             string `json:"type" valid:"required" enums:"audit,release"`                   // 节点类型，必填，选项为“审核”或“上线”
	InstanceName     string `json:"instance_name,omitempty" valid:"required_if=AuditType online"`  // 数据源名称，在线审核时必填
	InstanceType     string `json:"instance_type,omitempty" valid:"required_if=AuditType offline"` // 数据源类型，离线审核时必填
	ObjectPath       string `json:"object_path" valid:"required"`                                  // 审核脚本路径，必填，用户填写文件路径
	ObjectType       string `json:"object_type" valid:"required" enums:"sql,mybatis"`              // 审核对象类型，必填，可选项为SQL文件、MyBatis文件
	AuditMethod      string `json:"audit_method" valid:"required" enums:"offline,online"`          // 审核方式，必选，可选项为离线审核、在线审核
	RuleTemplateName string `json:"rule_template_name" valid:"required"`                           // 审核规则模板，必填
}

func (node *pipelineNodeBase) toSvcPipelineNode() *pipeline.PipelineNode {
	return &pipeline.PipelineNode{
		Name:             node.Name,
		NodeType:         node.Type,
		InstanceName:     node.InstanceName,
		InstanceType:     node.InstanceType,
		ObjectPath:       node.ObjectPath,
		ObjectType:       node.ObjectType,
		AuditMethod:      node.AuditMethod,
		RuleTemplateName: node.RuleTemplateName,
	}
}

// CreatePipelineReqV1 用于创建流水线的请求结构体
type CreatePipelineReqV1 struct {
	pipelineBase
	Nodes []pipelineNodeBase `json:"nodes" valid:"dive,required"` // 节点信息
}

func (req *CreatePipelineReqV1) convertToSvcPipeline(projectUID string) *pipeline.Pipeline {
	pipe := req.pipelineBase.toSvcPipeline()
	pipe.ProjectUID = projectUID
	pipe.PipelineNodes = make([]*pipeline.PipelineNode, 0, len(req.Nodes))
	for _, node := range req.Nodes {
		pipe.PipelineNodes = append(pipe.PipelineNodes, node.toSvcPipelineNode())
	}
	return pipe
}

type CreatePipelineResV1 struct {
	controller.BaseRes
	Data createPipelineResData `json:"data"`
}

type createPipelineResData struct {
	PipelineID uint `json:"pipeline_id"` // 流水线的唯一标识符
}

// @Summary 创建流水线
// @Description create pipeline
// @Id createPipelineV1
// @Tags pipeline
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param pipeline body v1.CreatePipelineReqV1 true "create pipeline"
// @Success 200 {object} v1.CreatePipelineResV1
// @router /v1/projects/{project_name}/pipelines [post]
func CreatePipeline(c echo.Context) error {
	// 1. 解析请求体中的数据
	req := new(CreatePipelineReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 2. 获取项目 UID
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	svcPipeline := req.convertToSvcPipeline(projectUid)

	var pipelineSvc pipeline.PipelineSvc

	// 3. 自动匹配与添加审核模板
	err = pipelineSvc.CheckRuleTemplate(svcPipeline)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = pipelineSvc.CheckInstance(c.Request().Context(), svcPipeline)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 4. 保存流水线
	err = pipelineSvc.CreatePipeline(svcPipeline, controller.GetUserID(c))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 5. 返回成功响应
	return c.JSON(http.StatusOK, &CreatePipelineResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: createPipelineResData{
			PipelineID: svcPipeline.ID,
		},
	})
}

// GetPipelinesReqV1 用于请求获取流水线列表的结构体
type GetPipelinesReqV1 struct {
	PageIndex           uint32 `json:"page_index" query:"page_index" valid:"required"`        // 页码索引
	PageSize            uint32 `json:"page_size" query:"page_size" valid:"required"`          // 每页条数
	FuzzySearchNameDesc string `json:"fuzzy_search_name_desc" query:"fuzzy_search_name_desc"` // 用于模糊搜索流水线名称和描述的关键字
}

// GetPipelinesResV1 用于响应流水线列表的结构体
type GetPipelinesResV1 struct {
	controller.BaseRes
	Data      []pipelineDetail `json:"data"`       // 流水线列表数据
	TotalNums uint64           `json:"total_nums"` // 流水线总数
}

// @Summary 获取流水线列表
// @Description get pipeline list
// @Id getPipelinesV1
// @Tags pipeline
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_name_desc query string false "fuzzy search pipeline name and description"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetPipelinesResV1
// @router /v1/projects/{project_name}/pipelines [get]
func GetPipelines(c echo.Context) error {
	// 1. 解析请求体中的数据
	req := new(GetPipelinesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 2. 获取项目 UID
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 3. 计算分页参数
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	// 4. 获取存储对象并查询流水线列表
	var pipelineSvc pipeline.PipelineSvc
	count, pipelineList, err := pipelineSvc.GetPipelineList(limit, offset, req.FuzzySearchNameDesc, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := make([]pipelineDetail, len(pipelineList))
	for idx, pipe := range pipelineList {
		data[idx].fillWith(pipe)
	}
	// 5. 返回成功响应
	return c.JSON(http.StatusOK, &GetPipelinesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		TotalNums: count,
		Data:      data,
	})
}

// GetPipelineDetailReqV1 用于请求获取流水线详情的结构体
type GetPipelineDetailReqV1 struct {
	PipelineID string `json:"pipeline_id" query:"pipeline_id" valid:"required"` // 流水线的唯一标识符
}

// GetPipelineDetailResV1 用于响应流水线详情的结构体
type GetPipelineDetailResV1 struct {
	controller.BaseRes
	Data pipelineDetailData `json:"data"`
}

// @Summary 获取流水线详情
// @Description get pipeline detail
// @Id getPipelineDetailV1
// @Tags pipeline
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param pipeline_id path string true "pipeline id"
// @Success 200 {object} v1.GetPipelineDetailResV1
// @router /v1/projects/{project_name}/pipelines/{pipeline_id}/ [get]
func GetPipelineDetail(c echo.Context) error {
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	pipelineID, err := strconv.Atoi(c.Param("pipeline_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var pipelineSvc pipeline.PipelineSvc
	pipe, err := pipelineSvc.GetPipeline(projectUid, uint(pipelineID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var pipelineDetail pipelineDetail
	pipelineDetail.fillWith(pipe)
	nodeDetails := make([]pipelineNodeDetail, len(pipe.PipelineNodes))
	for i, node := range pipe.PipelineNodes {
		nodeDetails[i].fillWith(c.Request().Context(), node, c.Param("project_name"))
	}

	return c.JSON(http.StatusOK, &GetPipelineDetailResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: pipelineDetailData{
			pipelineDetail: pipelineDetail,
			Nodes:          nodeDetails,
		},
	})
}

type pipelineDetailData struct {
	pipelineDetail
	Nodes []pipelineNodeDetail `json:"nodes"` // 流水线节点信息
}

// UpdatePipelineReqV1 用于更新流水线的请求结构体
type UpdatePipelineReqV1 struct {
	pipelineBase
	Nodes []updatePipelineNode `json:"nodes,omitempty" valid:"dive,required"` // 节点信息
}

type updatePipelineNode struct {
	ID uint `json:"id"`
	pipelineNodeBase
}

func (req UpdatePipelineReqV1) convertToSvcPipeline(projectUID string, pipelineID uint) *pipeline.Pipeline {
	svcPipeline := req.pipelineBase.toSvcPipeline()
	svcPipeline.ID = pipelineID
	svcPipeline.ProjectUID = projectUID
	svcPipeline.PipelineNodes = make([]*pipeline.PipelineNode, 0, len(req.Nodes))
	for idx, node := range req.Nodes {
		svcPipeline.PipelineNodes = append(svcPipeline.PipelineNodes, node.toSvcPipelineNode())
		svcPipeline.PipelineNodes[idx].ID = node.ID
	}
	return svcPipeline
}

// @Summary 更新流水线
// @Description update pipeline
// @Id updatePipelineV1
// @Tags pipeline
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param pipeline_id path string true "pipeline id"
// @Param pipeline body v1.UpdatePipelineReqV1 true "update pipeline"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/pipelines/{pipeline_id}/ [patch]
func UpdatePipeline(c echo.Context) error {

	req := new(UpdatePipelineReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	pipelineID, err := strconv.Atoi(c.Param("pipeline_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var pipelineSvc pipeline.PipelineSvc
	svcPipeline := req.convertToSvcPipeline(projectUid, uint(pipelineID))

	err = pipelineSvc.CheckRuleTemplate(svcPipeline)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = pipelineSvc.CheckInstance(c.Request().Context(), svcPipeline)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = pipelineSvc.UpdatePipeline(svcPipeline, controller.GetUserID(c))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// DeletePipelineReqV1 用于删除流水线的请求结构体
type DeletePipelineReqV1 struct {
	ProjectName string `json:"project_name" valid:"required"` // 项目名称，必填
	PipelineID  string `json:"pipeline_id" valid:"required"`  // 流水线 ID，必填
}

// @Summary 删除流水线
// @Description delete pipeline
// @Id deletePipelineV1
// @Tags pipeline
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param pipeline_id path string true "pipeline id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/pipelines/{pipeline_id}/ [delete]
func DeletePipeline(c echo.Context) error {
	// 获取项目的 UID
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	pipelineID, err := strconv.Atoi(c.Param("pipeline_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var pipelineSvc pipeline.PipelineSvc
	err = pipelineSvc.DeletePipeline(projectUid, uint(pipelineID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
