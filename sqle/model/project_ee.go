//go:build enterprise
// +build enterprise

package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

func (s *Storage) RemoveProjectRelateData(projectID ProjectUID) error {

	// 删除项目的频率不会很高, 对性能不敏感, 但删除的项很多, 一条SQL全删SQL会比较复杂, 所以将各功能模块分开删除, 以提高可维护性
	// 工单的删除复用回收工单的流程
	// 流程模板不用删除, 流程模板的绑定方式是 项目表中记录项目绑定的流程模板, 删除项目后流程模板自动作废
	return errors.ConnectStorageErrWrapper(s.Tx(func(txDB *gorm.DB) error {
		if err := s.deleteAllWhitelistByProjectID(txDB, projectID); err != nil {
			return err
		}

		if err := s.deleteAllRuleTemplateByProjectID(txDB, projectID); err != nil {
			return err
		}

		if err := s.deleteAllWorkflowByProjectID(txDB, projectID); err != nil {
			return err
		}

		// if err := s.deleteAllAuditPlanByProjectID(txDB, projectID); err != nil {
		// 	return err
		// }

		return nil
	}))
}

// 删除项目所有SQL白名单
func (s *Storage) deleteAllWhitelistByProjectID(tx *gorm.DB, projectID ProjectUID) error {
	return tx.Where("project_id = ?", projectID).Delete(&SqlWhitelist{}).Error
}

// 删除项目中所有规则模板
func (s *Storage) deleteAllRuleTemplateByProjectID(tx *gorm.DB, projectID ProjectUID) error {
	return tx.Where("project_id = ?", projectID).Delete(&RuleTemplate{}).Error
}

// 删除项目中所有工单/工单对应任务
func (s *Storage) deleteAllWorkflowByProjectID(tx *gorm.DB, projectID ProjectUID) error {
	workflows, err := s.GetWorkflowsByProjectID(projectID)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		err = s.deleteWorkflow(tx, workflow)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除项目中所有扫描任务
// func (s *Storage) deleteAllAuditPlanByProjectID(tx *gorm.DB, projectID ProjectUID) error {
// 	return tx.Where("project_id = ?", projectID).Delete(&AuditPlan{}).Error
// }

// // 删除项目中所有实例
// func (s *Storage) deleteAllInstanceByProjectID(tx *gorm.DB, projectID uint) error {
// 	return tx.Where("project_id = ?", projectID).Delete(&Instance{}).Error
// }

// // 删除项目本身
// func (s *Storage) deleteProjectByID(tx *gorm.DB, projectID uint) error {
// 	return tx.Where("id = ?", projectID).Delete(&Project{}).Error
// }

// func (s *Storage) GetProjectListBySyncTaskId(syncTaskID uint) ([]*Project, error) {
// 	projectList := make([]*Project, 0)
// 	err := s.db.Model(&Project{}).Preload("Instances", func(db *gorm.DB) *gorm.DB {
// 		return db.Where("sync_instance_task_id = ?", syncTaskID)
// 	}).Find(&projectList).Error
// 	if err != nil {
// 		return nil, errors.ConnectStorageErrWrapper(err)
// 	}

// 	var result []*Project
// 	for _, project := range projectList {
// 		if len(project.Instances) > 0 {
// 			result = append(result, project)
// 		}
// 	}

// 	return result, nil
// }

// func (s *Storage) ArchiveProject(projectName string) error {
// 	return s.db.Model(&Project{}).Where("name = ?", projectName).Update(Project{Status: ProjectStatusArchived}).Error
// }

// func (s *Storage) UnarchiveProject(projectName string) error {
// 	return s.db.Model(&Project{}).Where("name = ?", projectName).Update(Project{Status: ProjectStatusActive}).Error
// }
