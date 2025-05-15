package model

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

func init() {
	autoMigrateList = append(autoMigrateList, &Pipeline{})
	autoMigrateList = append(autoMigrateList, &PipelineNode{})
}

// 定义节点类型
type PipelineNodeType string

const (
	NodeTypeAudit   PipelineNodeType = "audit"
	NodeTypeRelease PipelineNodeType = "release"
)

// 定义审核对象类型
type ObjectType string

const (
	ObjectTypeSQL     ObjectType = "sql"
	ObjectTypeMyBatis ObjectType = "mybatis"
)

// 定义审核方式
type AuditMethod string

const (
	AuditMethodOffline AuditMethod = "offline"
	AuditMethodOnline  AuditMethod = "online"
)

type Pipeline struct {
	Model
	ProjectUid  ProjectUID `gorm:"index; not null" json:"project_uid"`     // 关联的流水线ID
	Name        string     `gorm:"type:varchar(255);not null" json:"name"` // 流水线名称
	Description string     `gorm:"type:varchar(512)" json:"description"`   // 流水线描述
	Address     string     `gorm:"type:varchar(255)" json:"address"`       // 关联流水线地址
}

type PipelineNode struct {
	gorm.Model
	PipelineID       uint   `gorm:"type:bigint;not null;index" json:"pipeline_id"`        // 关联的流水线ID
	UUID             string `gorm:"type:varchar(32);not null" json:"uuid"`                // 节点uuid
	Name             string `gorm:"type:varchar(255);not null" json:"name"`               // 节点名称
	NodeType         string `gorm:"type:varchar(20);not null" json:"node_type"`           // 节点类型
	NodeVersion      string `gorm:"type:varchar(32)" json:"node_version"`                 // 节点版本
	InstanceID       uint64 `gorm:"type:bigint" json:"instance_id"`                       // 数据源名称，在线审核时必填
	InstanceType     string `gorm:"type:varchar(255)" json:"instance_type,omitempty"`     // 数据源类型，离线审核时必填
	ObjectPath       string `gorm:"type:varchar(512);not null" json:"object_path"`        // 审核脚本路径
	ObjectType       string `gorm:"type:varchar(20);not null" json:"object_type"`         // 审核对象类型
	AuditMethod      string `gorm:"type:varchar(20);not null" json:"audit_method"`        // 审核方式
	RuleTemplateName string `gorm:"type:varchar(255);not null" json:"rule_template_name"` // 审核规则模板
	Token            string `gorm:"type:varchar(512);not null" json:"token"`              // token
}

func (p *PipelineNode) BeforeSave(tx *gorm.DB) (err error) {
	if !isValidPipelineNodeType(p.NodeType) {
		return fmt.Errorf("invalid node type: %s", p.NodeType)
	}
	if !isValidObjectType(p.ObjectType) {
		return fmt.Errorf("invalid object type: %s", p.ObjectType)
	}
	if !isValidAuditMethod(p.AuditMethod) {
		return fmt.Errorf("invalid audit method: %s", p.AuditMethod)
	}
	return nil
}

func isValidPipelineNodeType(t string) bool {
	for _, validType := range []PipelineNodeType{NodeTypeAudit, NodeTypeRelease} {
		if PipelineNodeType(t) == validType {
			return true
		}
	}
	return false
}

func isValidObjectType(o string) bool {
	for _, validObjectType := range []ObjectType{ObjectTypeSQL, ObjectTypeMyBatis} {
		if ObjectType(o) == validObjectType {
			return true
		}
	}
	return false
}

func isValidAuditMethod(a string) bool {
	for _, validMethod := range []AuditMethod{AuditMethodOffline, AuditMethodOnline} {
		if AuditMethod(a) == validMethod {
			return true
		}
	}
	return false
}

func (s *Storage) GetPipelineList(projectID ProjectUID, fuzzySearchContent string, limit, offset uint32) ([]*Pipeline, uint64, error) {
	var count int64
	var pipelines []*Pipeline
	query := s.db.Model(&Pipeline{}).Where("project_uid = ?", projectID)

	if fuzzySearchContent != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+fuzzySearchContent+"%", "%"+fuzzySearchContent+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return pipelines, uint64(count), errors.New(errors.ConnectStorageError, err)
	}

	if count == 0 {
		return pipelines, uint64(count), nil
	}

	err = query.Offset(int(offset)).Limit(int(limit)).Order("id desc").Find(&pipelines).Error
	return pipelines, uint64(count), errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetPipelineDetail(projectID ProjectUID, pipelineID uint) (*Pipeline, error) {
	pipeline := &Pipeline{}
	err := s.db.Model(Pipeline{}).Where("project_uid = ? AND id = ?", projectID, pipelineID).First(pipeline).Error
	if err != nil {
		return pipeline, errors.New(errors.ConnectStorageError, err)
	}
	return pipeline, nil
}

func (s *Storage) GetPipelineNode(pipelineID uint, nodeID uint) (*PipelineNode, error) {
	var node *PipelineNode
	err := s.db.Model(PipelineNode{}).Where("pipeline_id = ? AND id = ?", pipelineID, nodeID).First(&node).Error
	if err != nil {
		return node, errors.New(errors.ConnectStorageError, err)
	}
	return node, nil
}

func (s *Storage) GetPipelineNodes(pipelineID uint) ([]*PipelineNode, error) {
	var nodes []*PipelineNode
	err := s.db.Model(PipelineNode{}).Where("pipeline_id = ?", pipelineID).Find(&nodes).Error
	if err != nil {
		return nodes, errors.New(errors.ConnectStorageError, err)
	}
	return nodes, nil
}

func (s *Storage) GetPipelineNodesByInstanceId(instanceID uint64) ([]*PipelineNode, error) {
	if instanceID == 0 {
		return nil, fmt.Errorf("instance id should not be zero")
	}
	var nodes []*PipelineNode
	err := s.db.Model(PipelineNode{}).Where("instance_id = ?", instanceID).Find(&nodes).Error
	if err != nil {
		return nodes, errors.New(errors.ConnectStorageError, err)
	}
	return nodes, nil
}

func (s *Storage) CreatePipeline(pipeline *Pipeline, nodes []*PipelineNode) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 4.1 保存 Pipeline 到数据库
		if err := txDB.Create(pipeline).Error; err != nil {
			return fmt.Errorf("failed to create pipeline: %w", err)
		}
		// 4.2 创建 PipelineNodes 并保存到数据库
		for _, node := range nodes {
			node.PipelineID = pipeline.ID
			if err := txDB.Create(node).Error; err != nil {
				return fmt.Errorf("failed to create pipeline node: %w", err)
			}
		}
		return nil
	})
}

func (s *Storage) DeletePipeline(projectUID ProjectUID, pipelineID uint) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 删除 pipeline 相关的所有 nodes
		if err := txDB.Model(&PipelineNode{}).Where("pipeline_id = ?", pipelineID).Delete(&PipelineNode{}).Error; err != nil {
			return fmt.Errorf("failed to delete pipeline nodes: %w", err)
		}

		// 删除 pipeline
		if err := txDB.Model(&Pipeline{}).Where("project_uid = ? AND id = ?", projectUID, pipelineID).Delete(&Pipeline{}).Error; err != nil {
			return fmt.Errorf("failed to delete pipeline: %w", err)
		}

		return nil
	})
}

func (s *Storage) UpdatePipeline(pipe *Pipeline, newNodes []*PipelineNode) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 1 更新 pipeline
		err := txDB.Model(&Pipeline{}).Where("id = ? AND project_uid = ?", pipe.ID, pipe.ProjectUid).Updates(&pipe).Error
		if err != nil {
			return fmt.Errorf("failed to update pipeline: %w", err)
		}
		// 2 删除旧的 pipeline nodes
		if err := txDB.Where("pipeline_id = ?", pipe.ID).Delete(&PipelineNode{}).Error; err != nil {
			return fmt.Errorf("failed to delete old pipeline nodes: %w", err)
		}
		// 3 添加新的 pipeline nodes
		if err := txDB.CreateInBatches(&newNodes, 100).Error; err != nil {
			return fmt.Errorf("failed to update pipeline node: %w", err)
		}
		return nil
	})
}

func (s *Storage) UpdatePipelineNode(newNode *PipelineNode) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 3 更新节点属性
		if err := txDB.Save(&newNode).Error; err != nil {
			return fmt.Errorf("failed to update pipeline node: %w", err)
		}
		return nil
	})
}
