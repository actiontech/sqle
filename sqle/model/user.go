package model

type User struct {
	Model
	Name     string
	Mail     string
	Password string

	Roles []*Role `gorm:"many2many:user_role;"`
}

type Role struct {
	Model
	Name string
	Desc string

	Users     []*User     `gorm:"many2many:user_role;"`
	Instances []*Instance `gorm:"many2many:user_languages;"`
}
