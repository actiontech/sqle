package model

// import (
// 	"fmt"

// 	"github.com/actiontech/sqle/sqle/errors"

// 	"github.com/jinzhu/gorm"
// )

// type CloudBeaverUserCache struct {
// 	SQLEUserID        uint   `json:"sqle_user_id" gorm:"column:sqle_user_id;primary_key"`
// 	CloudBeaverUserID string `json:"cloud_beaver_user_id" gorm:"column:cloud_beaver_user_id"`
// 	SQLEFingerprint   string `json:"sqle_fingerprint" gorm:"column:sqle_fingerprint"`
// }

// func (s *Storage) GetCloudBeaverUserCacheByCBUserID(id string) (*CloudBeaverUserCache, bool, error) {
// 	c := &CloudBeaverUserCache{}
// 	err := s.db.Where("cloud_beaver_user_id = ?", id).Find(c).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return c, false, nil
// 	}

// 	return c, true, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) UpdateCloudBeaverUserCache(sqleUserID uint, CloudBeaverUserID string) error {
// 	user, exist, err := s.GetUserByID(sqleUserID)
// 	if err != nil {
// 		return err
// 	}
// 	if !exist {
// 		return errors.New(errors.ConnectStorageError, fmt.Errorf("sqle user not exist"))
// 	}

// 	err = s.db.Table("cloud_beaver_user_caches").Save(&CloudBeaverUserCache{
// 		SQLEUserID:        sqleUserID,
// 		CloudBeaverUserID: CloudBeaverUserID,
// 		SQLEFingerprint:   user.FingerPrint(),
// 	}).Error
// 	return errors.New(errors.ConnectStorageError, err)
// }

// type CloudBeaverInstanceCache struct {
// 	CloudBeaverInstanceID   string `json:"cloud_beaver_instance_id"`
// 	SQLEInstanceID          uint   `json:"sqle_instance_id" gorm:"primary_key"`
// 	SQLEInstanceFingerprint string `json:"sqle_instance_fingerprint"`
// }

// func (s *Storage) GetCloudBeaverInstanceCacheBySQLEInstIDs(ids []uint) ([]*CloudBeaverInstanceCache, error) {
// 	c := []*CloudBeaverInstanceCache{}
// 	err := s.db.Where("sqle_instance_id in (?)", ids).Find(&c).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return c, nil
// 	}

// 	return c, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) GetCloudBeaverInstanceCacheByCBInstIDs(ids []string) ([]*CloudBeaverInstanceCache, error) {
// 	c := []*CloudBeaverInstanceCache{}
// 	err := s.db.Where("cloud_beaver_instance_id in (?)", ids).Find(&c).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return c, nil
// 	}

// 	return c, errors.New(errors.ConnectStorageError, err)
// }
