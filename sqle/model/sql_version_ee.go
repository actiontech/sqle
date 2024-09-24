package model

import (
	"gorm.io/gorm"
)

func (s *Storage) GetNextSatgeByVersionIdAndSequence(txDB *gorm.DB, versionId uint, sequence int) (*SqlVersionStage, bool, error) {
	stage := &SqlVersionStage{}
	// next stage sequence
	next := sequence + 1
	err := txDB.Where("sql_version_id = ? AND stage_sequence = ?", versionId, next).First(stage).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return stage, true, nil
}
