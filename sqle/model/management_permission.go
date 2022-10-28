package model

// NOTE: related model:
// - model.User
type ManagementPermission struct {
	Model
	UserId         uint `gorm:"index"`
	PermissionCode uint `gorm:"comment:'平台管理权限'"`
}

const (
	// management permission list

	// 创建项目
	ManagementPermissionCreateProject = iota + 1
)
