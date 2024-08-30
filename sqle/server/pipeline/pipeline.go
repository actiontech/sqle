package pipeline

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	dmsCommonJwt "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/aliyun/credentials-go/credentials/utils"

	"gorm.io/gorm"
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

func (node PipelineNode) IntegrationInfo() string {
	dmsAddr := controller.GetDMSServerAddress()
	parsedURL, err := url.Parse(dmsAddr)
	if err != nil {
		return ""
	}
	ip, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return ""
	}

	switch model.PipelineNodeType(node.NodeType) {
	case model.NodeTypeAudit:
		var cmdUsage = "#使用方法#\n1. 确保运行该命令的用户具有scannerd的执行权限。\n2. 在scannerd文件所在目录执行启动命令。\n#启动命令#\n"
		baseCmd := "./scannerd %s --host=\"%s\" --port=\"%s\" --dir=\"%s\" --token=\"%s\""
		var extraArgs string
		var cmdType string
		if model.ObjectType(node.ObjectType) == model.ObjectTypeSQL {
			cmdType = "sql_file"
		}
		if model.ObjectType(node.ObjectType) == model.ObjectTypeMyBatis {
			cmdType = "mysql_mybatis"
		}
		if model.AuditMethod(node.AuditMethod) == model.AuditMethodOnline {
			extraArgs = fmt.Sprintf(" --instance-name=\"%s\"", node.InstanceName)
		}
		if model.AuditMethod(node.AuditMethod) == model.AuditMethodOffline {
			extraArgs = fmt.Sprintf(" --db-type=\"%s\"", node.InstanceType)
		}
		return fmt.Sprintf(cmdUsage+baseCmd+extraArgs, cmdType, ip, port, node.ObjectPath, node.Token)
	case model.NodeTypeRelease:
		return ""
	default:
		return ""
	}
}

type PipelineNode struct {
	ID               uint   // 节点的唯一标识符，在更新时必填
	Version          string // 节点版本ID
	Name             string // 节点名称，必填，支持中文、英文+数字+特殊字符
	NodeType         string // 节点类型，必填，选项为“审核”或“上线”
	InstanceName     string // 数据源名称，在线审核时必填
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
		}
	}
	return nil
}

func (svc PipelineSvc) CreatePipeline(pipe *Pipeline, userID string) error {
	s := model.GetStorage()
	modelPipeline := svc.toModelPipeline(pipe)
	modelPipelineNodes := svc.toModelPipelineNodes(pipe, userID)
	err := s.CreatePipeline(modelPipeline, modelPipelineNodes)
	if err != nil {
		return err
	}
	pipe.ID = modelPipeline.ID
	return nil
}

func (svc PipelineSvc) toModelPipeline(pipe *Pipeline) *model.Pipeline {
	if pipe == nil {
		return nil
	}

	return &model.Pipeline{
		ProjectUid:  model.ProjectUID(pipe.ProjectUID),
		Name:        pipe.Name,
		Description: pipe.Description,
		Address:     pipe.Address,
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
			InstanceName:     node.InstanceName,
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

func (svc PipelineSvc) GetPipelineList(limit, offset uint32, fuzzySearchNameDesc string, projectUID string) (count uint64, pipelines []*Pipeline, err error) {
	s := model.GetStorage()
	modelPipelines, count, err := s.GetPipelineList(model.ProjectUID(projectUID), fuzzySearchNameDesc, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	pipelines = make([]*Pipeline, 0, len(modelPipelines))
	for _, modelPipeline := range modelPipelines {
		modelPiplineNodes, err := s.GetPipelineNodes(modelPipeline.ID)
		if err != nil {
			return 0, nil, err
		}
		pipelines = append(pipelines, svc.toPipeline(modelPipeline, modelPiplineNodes))
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
		InstanceName:     modelPipelineNode.InstanceName,
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
		newNode.InstanceName != oldNode.InstanceName ||
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
	4. uv2表示用户理解的版本2，此时由于node1实际未修改，因此继承了之前节点的uuid version token，用户可以无需修改命令
	5. 在用户修改的第三个版本中，即行uv3，节点1更新了参数(例如修改了审核路径，实际上用户也会知道需要修改命令)，导致需要修改命令，因此这里的uuid version token都更新了
	版本
	↓
	+-------+-------+-------+-------+-------+-------+
	|       |     node1     |     node2     | node3 | <- 节点名称
	+-------+-------+-------+-------+-------+-------+
	|       | uuid1 | uuid2 | uuid3 | uuid4 | uuid5 | <- 此列不展示
	+-------+-------+-------+-------+-------+-------+
	|  uv1  | n1-v1 |       | n2-v1 |       |       |
	+-------+-------+-------+-------+-------+-------+
	|  uv2  | n1-v1 |       |       | n2-v2 |       |
	+-------+-------+-------+-------+-------+-------+
	|  uv3  |       | n1-v3 |       | n2-v2 |       |
	+-------+-------+-------+-------+-------+-------+
	|  uv4  |       | n1-v3 |       | n2-v2 | n3-v4 |
	+-------+-------+-------+-------+-------+-------+

```
*/
func (svc PipelineSvc) UpdatePipeline(pipe *Pipeline, userId string) error {
	s := model.GetStorage()
	return s.Tx(func(txDB *gorm.DB) error {
		// 4.1 更新 pipeline
		err := txDB.Model(&model.Pipeline{}).
			Where("id = ? AND project_uid = ?", pipe.ID, pipe.ProjectUID).
			Updates(model.Pipeline{
				Name:        pipe.Name,
				Description: pipe.Description,
				Address:     pipe.Address,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to update pipeline: %w", err)
		}
		nodes := []model.PipelineNode{}
		if err := txDB.Where("pipeline_id = ?", pipe.ID).Find(&nodes).Error; err != nil {
			return fmt.Errorf("failed to delete old pipeline nodes: %w", err)
		}
		oldNodeMap := make(map[uint] /* node id */ *model.PipelineNode, len(nodes))
		for idx, node := range nodes {
			oldNodeMap[node.ID] = &nodes[idx]
		}
		// 4.2 删除旧的 pipeline nodes
		if err := txDB.Where("pipeline_id = ?", pipe.ID).Delete(&model.PipelineNode{}).Error; err != nil {
			return fmt.Errorf("failed to delete old pipeline nodes: %w", err)
		}
		v := svc.newVersion()
		// 4.3 添加新的 pipeline nodes
		var newToken string
		var newUuid string
		var newVersion string
		for _, newNode := range pipe.PipelineNodes {
			/*
				若节点命令不变更，则节点继承原有节点对应的token uuid和version
				若节点命令变更，则节点使用新的token uuid和version
			*/
			oldNode, exist := oldNodeMap[newNode.ID]
			if exist {
				newToken = oldNode.Token
				newUuid = oldNode.UUID
				newVersion = oldNode.NodeVersion
			}
			if !exist || svc.needUpdateToken(oldNode, newNode) {
				newVersion = v
				newUuid = utils.GetUUID()
				newToken, err = svc.newToken(userId, newVersion, newUuid)
				if err != nil {
					return err
				}
			}
			node := model.PipelineNode{
				PipelineID:       pipe.ID,
				Name:             newNode.Name,
				NodeType:         newNode.NodeType,
				InstanceName:     newNode.InstanceName,
				InstanceType:     newNode.InstanceType,
				ObjectPath:       newNode.ObjectPath,
				ObjectType:       newNode.ObjectType,
				AuditMethod:      newNode.AuditMethod,
				RuleTemplateName: newNode.RuleTemplateName,
				NodeVersion:      newVersion,
				Token:            newToken,
				UUID:             newUuid,
			}

			if err := txDB.Create(&node).Error; err != nil {
				return fmt.Errorf("failed to update pipeline node: %w", err)
			}
		}
		return nil
	})
}

func (svc PipelineSvc) DeletePipeline(projectUID string, pipelineID uint) error {
	s := model.GetStorage()
	return s.Tx(func(txDB *gorm.DB) error {
		// 删除 pipeline 相关的所有 nodes
		if err := txDB.Model(&model.PipelineNode{}).Where("pipeline_id = ?", pipelineID).Delete(&model.PipelineNode{}).Error; err != nil {
			return fmt.Errorf("failed to delete pipeline nodes: %w", err)
		}

		// 删除 pipeline
		if err := txDB.Model(&model.Pipeline{}).Where("id = ?", pipelineID).Delete(&model.Pipeline{}).Error; err != nil {
			return fmt.Errorf("failed to delete pipeline: %w", err)
		}

		return nil
	})
}
