package pipeline

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsCommonJwt "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"
	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/aliyun/credentials-go/credentials/utils"
)

type Pipeline struct {
	ID            uint            // 流水线的唯一标识符
	ProjectUID    string          // 项目UID
	Token         string          // token
	Name          string          // 流水线名称
	Description   string          // 流水线描述
	Address       string          // 关联流水线地址
	PipelineNodes []*PipelineNode // 节点
}

func (pipe Pipeline) NodeCount() uint32 {
	return uint32(len(pipe.PipelineNodes))
}

func (node PipelineNode) IntegrationInfo(ctx context.Context, projectName string) (string, error) {
	dmsAddr := controller.GetDMSServerAddress()
	parsedURL, err := url.Parse(dmsAddr)
	if err != nil {
		return "", err
	}
	ip, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return "", err
	}
	if node.InstanceID != 0 {
		instance, exist, err := dms.GetInstancesById(context.TODO(), fmt.Sprint(node.InstanceID))
		if err != nil {
			return "", err
		}
		if !exist {
			return "", errors.NewInstanceNoExistErr()
		}
		node.InstanceName = instance.Name
	}

	switch model.PipelineNodeType(node.NodeType) {
	case model.NodeTypeAudit:
		cmdUsage := locale.Bundle.LocalizeMsgByCtx(ctx, locale.PipelineCmdUsage)

		var cmd string
		var cmdType string
		if model.ObjectType(node.ObjectType) == model.ObjectTypeSQL {
			cmdType = scannerCmd.TypeSQLFile
		}
		if model.ObjectType(node.ObjectType) == model.ObjectTypeMyBatis {
			cmdType = scannerCmd.TypeMySQLMybatis
		}
		sqlfile, err := scannerCmd.GetScannerdCmd(cmdType)
		if err != nil {
			return "", err
		}
		params := map[string]string{
			scannerCmd.FlagHost:      ip,
			scannerCmd.FlagPort:      port,
			scannerCmd.FlagToken:     node.Token,
			scannerCmd.FlagDirectory: node.ObjectPath,
			scannerCmd.FlagDbType:    node.InstanceType,
			scannerCmd.FlagProject:   projectName,
		}
		if node.InstanceName != "" {
			params[scannerCmd.FlagInstanceName] = node.InstanceName
		}
		cmd, err = sqlfile.GenCommand("./scannerd", params)
		if err != nil {
			return "", err
		}
		return cmdUsage + cmd, nil
	case model.NodeTypeRelease:
		return "", fmt.Errorf("unsupport node type release")
	default:
		return "", fmt.Errorf("unsupport node type unknown")
	}
}

type PipelineNode struct {
	ID               uint   // 节点的唯一标识符，在更新时必填
	Version          string // 节点版本ID
	Name             string // 节点名称，必填，支持中文、英文+数字+特殊字符
	NodeType         string // 节点类型，必填，选项为“审核”或“上线”
	InstanceName     string // 数据源名称，在线审核时必填
	InstanceID       uint64 // 数据源ID
	InstanceType     string // 数据源类型，在线审核时必填
	ObjectPath       string // 审核脚本路径，必填，用户填写文件路径
	ObjectType       string // 审核对象类型，必填，可选项为SQL文件、MyBatis文件
	AuditMethod      string // 审核方式，必选，可选项为离线审核、在线审核
	RuleTemplateName string // 审核规则模板，必填
	Token            string // 节点Token
}

type PipelineSvc struct{}

func (svc PipelineSvc) CheckRuleTemplate(pipe *Pipeline) (err error) {
	s := model.GetStorage()
	for _, node := range pipe.PipelineNodes {
		exist, err := s.IsRuleTemplateExist(node.RuleTemplateName, []string{
			pipe.ProjectUID,
			model.ProjectIdForGlobalRuleTemplate,
		})
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("rule template does not exist")
		}
	}
	return nil
}

func (svc PipelineSvc) CheckInstance(ctx context.Context, pipe *Pipeline) (err error) {
	for _, node := range pipe.PipelineNodes {
		if node.InstanceName != "" {
			instance, exist, err := dms.GetInstanceInProjectByName(ctx, pipe.ProjectUID, node.InstanceName)
			if err != nil {
				return err
			}
			if !exist {
				return fmt.Errorf("instance does not exist")
			}
			node.InstanceType = instance.DbType
			node.InstanceID = instance.ID
		}
	}
	return nil
}

func (svc PipelineSvc) CreatePipeline(pipe *Pipeline, userID string) error {
	s := model.GetStorage()
	modelPipeline := svc.toModelPipeline(pipe, userID)
	modelPipelineNodes := svc.toModelPipelineNodes(pipe, userID)
	err := s.CreatePipeline(modelPipeline, modelPipelineNodes)
	if err != nil {
		return err
	}
	pipe.ID = modelPipeline.ID
	return nil
}

func (svc PipelineSvc) toModelPipeline(pipe *Pipeline, userId string) *model.Pipeline {
	if pipe == nil {
		return nil
	}

	return &model.Pipeline{
		ProjectUid:   model.ProjectUID(pipe.ProjectUID),
		Name:         pipe.Name,
		Description:  pipe.Description,
		Address:      pipe.Address,
		CreateUserID: userId,
	}
}

func (svc PipelineSvc) toModelPipelineNodes(pipe *Pipeline, userId string) []*model.PipelineNode {
	if pipe == nil || len(pipe.PipelineNodes) == 0 {
		return nil
	}
	nodeVersion := svc.newVersion()
	modelNodes := make([]*model.PipelineNode, 0, len(pipe.PipelineNodes))
	for _, node := range pipe.PipelineNodes {
		nodeUuid := utils.GetUUID()
		token, err := svc.newToken(userId, nodeVersion, nodeUuid)
		if err != nil {
			return nil
		}
		modelNode := &model.PipelineNode{
			PipelineID:       pipe.ID, // 需要将 Pipeline 的 ID 关联到 Node 上
			Name:             node.Name,
			NodeType:         node.NodeType,
			InstanceID:       node.InstanceID,
			InstanceType:     node.InstanceType,
			ObjectPath:       node.ObjectPath,
			ObjectType:       node.ObjectType,
			AuditMethod:      node.AuditMethod,
			RuleTemplateName: node.RuleTemplateName,
			NodeVersion:      nodeVersion,
			UUID:             nodeUuid,
			Token:            token,
		}
		modelNodes = append(modelNodes, modelNode)
	}

	return modelNodes
}

func (svc PipelineSvc) newVersion() string {
	return utils.GetUUID()
}

/*
1. token 由用户uid、node版本id、node版本下的uid和过期时间信息构成
2. 因此同一版本每一个node都有一个唯一的token, token能够通过版本号+时间戳索引到唯一的node
3. 不同版本的node，当配置参数不变，可以通过继承node上一版本的token和版本下的uid和token，达到不需要启动命令的效果
*/
func (svc PipelineSvc) newToken(userId, version, uuid string) (string, error) {
	token, err := dmsCommonJwt.GenJwtToken(
		dmsCommonJwt.WithUserId(userId),
		dmsCommonJwt.WithExpiredTime(365*24*time.Hour),
		// TODO 带上版本和uuid信息
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (svc PipelineSvc) GetPipeline(projectUID string, pipelineID uint) (*Pipeline, error) {
	s := model.GetStorage()
	modelPipeline, err := s.GetPipelineDetail(model.ProjectUID(projectUID), pipelineID)
	if err != nil {
		return nil, err
	}
	modelPiplineNodes, err := s.GetPipelineNodes(pipelineID)
	if err != nil {
		return nil, err
	}
	return svc.toPipeline(modelPipeline, modelPiplineNodes), nil
}

// GetPipelineListWithPermission 根据用户权限获取流水线列表
func (svc PipelineSvc) GetPipelineListWithPermission(limit, offset uint32, fuzzySearchNameDesc string, projectUID string, userPermission *dms.UserPermission, userId string) (count uint64, pipelines []*Pipeline, err error) {
	s := model.GetStorage()

	// 根据用户权限确定查询参数
	var queryUserId string
	var rangeDatasourceIds []string
	var canViewAll bool

	// 权限判断逻辑
	if userPermission.IsAdmin() || userPermission.IsProjectAdmin() {
		// 超级管理员或项目管理员：可以查看所有流水线
		canViewAll = true
	} else if viewPipelinePermission := userPermission.GetOnePermission(v1.OpPermissionViewPipeline); viewPipelinePermission != nil {
		// 拥有"查看流水线"权限的普通用户：可以查看指定数据源相关的流水线 + 自己创建的所有流水线
		queryUserId = userId
		rangeDatasourceIds = viewPipelinePermission.RangeUids
		canViewAll = false
	} else {
		// 普通用户：只能查看自己创建的流水线
		queryUserId = userId
		rangeDatasourceIds = nil
		canViewAll = false
	}

	// 执行数据库查询
	modelPipelines, count, err := s.GetPipelineList(model.ProjectUID(projectUID), fuzzySearchNameDesc, limit, offset, queryUserId, rangeDatasourceIds, canViewAll)
	if err != nil {
		return 0, nil, err
	}

	// 转换为服务层对象
	pipelines = make([]*Pipeline, 0, len(modelPipelines))
	if len(modelPipelines) == 0 {
		return count, pipelines, nil
	}

	// 收集所有pipeline ID
	pipelineIDs := make([]uint, 0, len(modelPipelines))
	for _, mp := range modelPipelines {
		pipelineIDs = append(pipelineIDs, mp.ID)
	}

	// 批量获取所有节点
	nodesMap, err := s.GetPipelineNodesInBatch(pipelineIDs)
	if err != nil {
		return 0, nil, err
	}

	// 组装结果
	for _, modelPipeline := range modelPipelines {
		nodes := nodesMap[modelPipeline.ID]
		pipelines = append(pipelines, svc.toPipeline(modelPipeline, nodes))
	}
	return count, pipelines, nil
}

func (svc PipelineSvc) toPipeline(modelPipeline *model.Pipeline, modelPipelineNodes []*model.PipelineNode) *Pipeline {
	if modelPipeline == nil {
		return nil
	}
	pipeline := &Pipeline{
		ID:          modelPipeline.ID,
		ProjectUID:  string(modelPipeline.ProjectUid), // 如果 ProjectUID 是字符串，可以直接转换
		Name:        modelPipeline.Name,
		Description: modelPipeline.Description,
		Address:     modelPipeline.Address,
	}
	if len(modelPipelineNodes) > 0 {
		pipeline.PipelineNodes = make([]*PipelineNode, 0, len(modelPipelineNodes))
		for _, node := range modelPipelineNodes {
			pipeline.PipelineNodes = append(pipeline.PipelineNodes, svc.toPipelineNode(node))
		}
	}
	return pipeline
}

func (svc PipelineSvc) toPipelineNode(modelPipelineNode *model.PipelineNode) *PipelineNode {
	if modelPipelineNode == nil {
		return nil
	}
	return &PipelineNode{
		ID:               modelPipelineNode.ID,
		Name:             modelPipelineNode.Name,
		NodeType:         modelPipelineNode.NodeType,
		InstanceID:       modelPipelineNode.InstanceID,
		InstanceType:     modelPipelineNode.InstanceType,
		ObjectPath:       modelPipelineNode.ObjectPath,
		ObjectType:       modelPipelineNode.ObjectType,
		AuditMethod:      modelPipelineNode.AuditMethod,
		RuleTemplateName: modelPipelineNode.RuleTemplateName,
		Token:            modelPipelineNode.Token,
	}
}

/*
needUpdateToken比较旧节点（oldNode）和新节点（newNode）的属性

	若以下任意属性发生变化，则认为需要更新 token
		1. RuleTemplateName：规则模板名称
		2. NodeType：节点类型
		3. ObjectPath：对象路径
		4. ObjectType：对象类型
		5. AuditMethod：审计方法
		6. InstanceName：实例名称
		7. InstanceType：实例类型
*/
func (svc PipelineSvc) needUpdateToken(oldNode *model.PipelineNode, newNode *PipelineNode) bool {
	return newNode.RuleTemplateName != oldNode.RuleTemplateName ||
		newNode.NodeType != oldNode.NodeType ||
		newNode.ObjectPath != oldNode.ObjectPath ||
		newNode.ObjectType != oldNode.ObjectType ||
		newNode.AuditMethod != oldNode.AuditMethod ||
		newNode.InstanceID != oldNode.InstanceID ||
		newNode.InstanceType != oldNode.InstanceType
}

/*
	*节点更新的逻辑有以下功能需求：

	1. 流水线节点更新一次，对用户来说视为一个版本，数据存储上需要能够方便筛选出每一个版本的所有节点
	2. 流水线的每一个节点，会生成一段执行脚本的命令，当节点信息修改可以不影响脚本命令时，不应改变脚本
	3. Scannerd在客户的流水线中执行，请求SQLE时携带的token需要能够唯一确定SQLE上的一个流水线节点

	*节点更新时，遵循以下逻辑：

	1. 当用户没有修改会导致需要变更启动命令的配置时，节点version uuid token继承之前节点的数据
	2. 当用户修改了节点配置，导致该节点要正确运行需要重新配置运行命令时，节点会使用新的version uuid token

```

	*在界面上的表示形态如下图所示：
	1. uv1表示的是用户理解的版本1 user version 1
	2. node1 表示用户理解的节点1，uuid表示节点的唯一id
	3. n1-v1表示节点1的第一个版本 node 1 version 1
	4. uv2表示用户理解的版本2，此时由于node1实际未修改，因此继承了之前节点的uuid version token，用户可以无需修改启动命令
	5. 在用户修改的第三个版本中，即行uv3，节点1更新了参数(例如修改了审核路径，实际上用户也会知道需要修改命令)，导致需要修改启动命令，因此这里的version token都更新了
	版本
	↓
	+-------+-------+-------+-------+
	|       | node1 | node2 | node3 | <- 节点名称
	+-------+-------+-------+-------+
	|       | uuid1 | uuid2 | uuid3 | <- 此列不展示
	+-------+-------+-------+-------+
	|  uv1  | n1-v1 | n2-v1 |       |
	+-------+-------+-------+-------+
	|  uv2  | n1-v1 | n2-v2 |       |
	+-------+-------+-------+-------+
	|  uv3  | n1-v3 | n2-v2 |       |
	+-------+-------+-------+-------+
	|  uv4  | n1-v3 | n2-v2 | n3-v4 |
	+-------+-------+-------+-------+

```
*/
func (svc PipelineSvc) UpdatePipeline(pipe *Pipeline, userId string) error {
	s := model.GetStorage()
	currentNodes, err := s.GetPipelineNodes(pipe.ID)
	if err != nil {
		return err
	}
	oldNodeMap := make(map[uint] /* node id */ *model.PipelineNode, len(currentNodes))
	for idx, node := range currentNodes {
		oldNodeMap[node.ID] = currentNodes[idx]
	}
	v := svc.newVersion()
	var newToken string
	var newUuid string
	var newVersion string
	newNodes := make([]*model.PipelineNode, 0, len(pipe.PipelineNodes))
	for _, newNode := range pipe.PipelineNodes {
		/*
			若节点命令不变更，则节点继承原有节点对应的token uuid和version
			若节点命令变更，则节点使用新的token uuid和version
		*/
		// 节点不变
		oldNode, exist := oldNodeMap[newNode.ID]
		if exist {
			newToken = oldNode.Token
			newUuid = oldNode.UUID
			newVersion = oldNode.NodeVersion
			// 当节点为在线审核，没有选择数据源时，当且仅当在更新时，用户没有修改数据源，才会出现，此时默认为原来的数据源id
			if newNode.AuditMethod == string(model.AuditMethodOnline) && newNode.InstanceID == 0 {
				newNode.InstanceID = oldNode.InstanceID
			}
		}
		// 更新节点
		if exist && svc.needUpdateToken(oldNode, newNode) {
			newVersion = v
			newToken, err = svc.newToken(userId, newVersion, newUuid)
			if err != nil {
				return err
			}
		}
		// 新建节点
		if !exist {
			newVersion = v
			newUuid = utils.GetUUID()
			newToken, err = svc.newToken(userId, newVersion, newUuid)
			if err != nil {
				return err
			}
		}
		newNodes = append(newNodes, &model.PipelineNode{
			PipelineID:       pipe.ID,
			Name:             newNode.Name,
			NodeType:         newNode.NodeType,
			InstanceID:       newNode.InstanceID,
			InstanceType:     newNode.InstanceType,
			ObjectPath:       newNode.ObjectPath,
			ObjectType:       newNode.ObjectType,
			AuditMethod:      newNode.AuditMethod,
			RuleTemplateName: newNode.RuleTemplateName,
			NodeVersion:      newVersion,
			Token:            newToken,
			UUID:             newUuid,
		})
	}
	newPipe := &model.Pipeline{
		Model: model.Model{
			ID: pipe.ID,
		},
		ProjectUid:  model.ProjectUID(pipe.ProjectUID),
		Name:        pipe.Name,
		Description: pipe.Description,
		Address:     pipe.Address,
	}
	return s.UpdatePipeline(newPipe, newNodes)
}

func (svc PipelineSvc) DeletePipeline(projectUID string, pipelineID uint) error {
	s := model.GetStorage()
	return s.DeletePipeline(model.ProjectUID(projectUID), pipelineID)
}

func (svc PipelineSvc) RefreshPipelineToken(pipelineID uint, nodeID uint, userID string) error {
	s := model.GetStorage()
	modelPiplineNode, err := s.GetPipelineNode(pipelineID, nodeID)
	if err != nil {
		return err
	}
	newToken, err := svc.newToken(userID, modelPiplineNode.NodeVersion, modelPiplineNode.UUID)
	if err != nil {
		return err
	}
	modelPiplineNode.Token = newToken
	return s.UpdatePipelineNode(modelPiplineNode)
}
