package model

// import (
// 	"github.com/actiontech/sqle/sqle/errors"
// 	"github.com/actiontech/sqle/sqle/utils"

// 	"github.com/jinzhu/gorm"
// )

// // NOTE: related model:
// // - model.User
// type ManagementPermission struct {
// 	Model
// 	UserId         uint `gorm:"index"`
// 	PermissionCode uint `gorm:"comment:'平台管理权限'"`
// }

// const (
// 	// management permission list

// 	// 创建项目
// 	ManagementPermissionCreateProject uint = iota + 1
// )

// var managementPermission2Desc = map[uint]string{
// 	ManagementPermissionCreateProject: "创建项目",
// }

// func GetManagementPermissionDesc(code uint) string {
// 	desc, ok := managementPermission2Desc[code]
// 	if !ok {
// 		desc = "未知权限"
// 	}
// 	return desc
// }

// func GetManagementPermission() map[uint] /*code*/ string /*desc*/ {
// 	resp := map[uint]string{}
// 	for u, s := range managementPermission2Desc {
// 		resp[u] = s
// 	}
// 	return resp
// }

// func (s *Storage) UpdateManagementPermission(userID uint, permissionCode []uint) error {
// 	return s.Tx(func(txDB *gorm.DB) error {
// 		err := updateManagementPermission(txDB, userID, permissionCode)
// 		if err != nil {
// 			txDB.Rollback()
// 			return errors.ConnectStorageErrWrapper(err)
// 		}
// 		return nil
// 	})
// }

// func updateManagementPermission(txDB *gorm.DB, userID uint, permissionCode []uint) error {
// 	if err := txDB.Where("user_id = ?", userID).Delete(&ManagementPermission{}).Error; err != nil {
// 		return err
// 	}
// 	for _, code := range permissionCode {
// 		if err := txDB.Create(&ManagementPermission{
// 			UserId:         userID,
// 			PermissionCode: code,
// 		}).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *Storage) GetManagementPermissionByUserID(userID uint) ([]uint, error) {
// 	code := []*struct {
// 		PermissionCode uint
// 	}{}
// 	err := s.db.Table("management_permissions").Select("permission_code").Where("user_id = ?", userID).Where("deleted_at IS NULL").Find(&code).Error
// 	resp := []uint{}
// 	for _, c := range code {
// 		resp = append(resp, c.PermissionCode)
// 	}
// 	return resp, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) GetManagementPermissionByUserIDs(userIDs []uint) (map[uint] /*user id*/ []uint /*codes*/, error) {
// 	p := []*ManagementPermission{}
// 	err := s.db.Table("management_permissions").Select("user_id,permission_code").Where("user_id in (?)", userIDs).Where("deleted_at IS NULL").Scan(&p).Error

// 	resp := map[uint][]uint{}
// 	for _, permission := range p {
// 		resp[permission.UserId] = append(resp[permission.UserId], permission.PermissionCode)
// 	}

// 	for id := range resp {
// 		resp[id] = utils.RemoveDuplicateUint(resp[id])
// 	}

// 	return resp, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) CheckUserHaveManagementPermission(userID uint, code []uint) (bool, error) {
// 	code = utils.RemoveDuplicateUint(code)

// 	user, _, err := s.GetUserByID(userID)
// 	if err != nil {
// 		return false, err
// 	}
// 	if user.Name == DefaultAdminUser {
// 		return true, nil
// 	}

// 	var count int
// 	err = s.db.Model(&ManagementPermission{}).
// 		Where("user_id = ?", userID).
// 		Where("permission_code in (?)", code).
// 		Count(&count).Error

// 	return count == len(code), errors.ConnectStorageErrWrapper(err)
// }
