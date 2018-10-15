package storage

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string
	Password string
	Roles    []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	gorm.Model
	Name       string
	Privileges []Privilege `gorm:"many2many:role_privileges"`
}

type Privilege struct {
	gorm.Model
	Code int
}

// privilege code
const (
	PRIVILEGE_CREATE_TASK = iota
)

func (s *Storage) GetUserById(id string) (*User, error) {
	user := &User{}
	err := s.db.First(user, id).Error
	return user, err
}

func (s *Storage) GetUserByName(name string) (*User, error) {
	user := &User{}
	err := s.db.Where("name = ?", name).First(user).Error
	return user, err
}

func (s *Storage) UpdateUser(user *User) error {
	return s.db.Save(user).Error
}

func (s *Storage) DelUser(user *User) error {
	return s.db.Delete(user).Error
}

func (s *Storage) GetUsers() ([]*User, error) {
	users := []*User{}
	err := s.db.Find(users).Error
	return users, err
}
