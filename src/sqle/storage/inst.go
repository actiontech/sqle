package storage

import "github.com/jinzhu/gorm"

const (
	DB_TYPE_MYSQL = iota
	DB_TYPE_MYCAT
	DB_TYPE_SQLSERVER
)

// Instance is a table for database info
type Instance struct {
	Model
	Name          string         `json:"name" gorm:"not null;index" example:""`
	DbType        int            `json:"type" gorm:"not null" example:"0"`
	Host          string         `json:"host" gorm:"not null" example:"10.10.10.10"`
	Port          string         `json:"port" gorm:"not null" example:"3306"`
	User          string         `json:"user" gorm:"not null" example:"root"`
	Password      string         `json:"-" gorm:"not null"`
	Desc          string         `json:"desc" example:"this is a instance"`
	RuleTemplates []RuleTemplate `json:"-" gorm:"many2many:instance_rule_template"`
}

func (s *Storage) GetInstById(id string) (*Instance, bool, error) {
	inst := &Instance{}
	err := s.db.Where("id = ?", id).First(inst).Error
	if err == gorm.ErrRecordNotFound {
		return inst, false, nil
	}
	return inst, true, err
}

func (s *Storage) GetInstByName(name string) (*Instance, bool, error) {
	inst := &Instance{}
	err := s.db.Where("name = ?", name).First(inst).Error
	if err == gorm.ErrRecordNotFound {
		return inst, false, nil
	}
	return inst, true, err
}

func (s *Storage) UpdateInst(inst *Instance) error {
	return s.db.Save(inst).Error
}

func (s *Storage) DelInstByName(inst *Instance) error {
	return s.db.Delete(inst).Error
}

func (s *Storage) GetDatabases() ([]Instance, error) {
	inst := []Instance{}
	err := s.db.Find(&inst).Error
	return inst, err
}

func (s *Storage) UpdateInstRuleTemplate(inst *Instance, ts ...RuleTemplate) error {
	return s.db.Model(inst).Association("RuleTemplates").Append(ts).Error
}
