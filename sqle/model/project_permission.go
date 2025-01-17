package model

import "github.com/actiontech/sqle/sqle/errors"

func (s *Storage) UserCanAccessWorkflow(userId string, workflowId string) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.workflow_id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
where w.deleted_at IS NULL
AND (w.create_user_id = ? OR cur_ws.assignees REGEXP ?)
`
	var count int64
	err := s.db.Raw(query, workflowId, userId, userId).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

func (s *Storage) UserCanViewWorkflow(userId string, workflowId string) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.workflow_id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_steps AS op_ws ON w.workflow_id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
where w.deleted_at IS NULL
AND (w.create_user_id = ? OR cur_ws.assignees REGEXP ? OR op_ws.operation_user_id = ?)
`
	var count int64
	err := s.db.Raw(query, workflowId, userId, userId, userId).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}
