//go:build !enterprise
// +build !enterprise

package model

func (s *Storage) GetOperationRecordProjectNameList() ([]string, error) {
	var projectNameList []string
	err := s.db.Model(&OperationRecord{}).Group("operation_project_name").Pluck("operation_project_name", &projectNameList).Error
	if err != nil {
		return nil, err
	}
	return projectNameList, err
}
