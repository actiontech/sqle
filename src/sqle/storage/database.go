package storage

import "github.com/jinzhu/gorm"

const (
	DB_TYPE_MYSQL = iota
	DB_TYPE_MYCAT
	DB_TYPE_SQLSERVER
)

type Db struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" gorm:"unique;not null;unique_index"`
	DbType     int    `json:"type" gorm:"not null"`
	Host       string `json:"host" gorm:"not null"`
	Port       string `json:"port" gorm:"not null"`
	User       string `json:"user" gorm:"not null"`
	Password   string `json:"-" gorm:"not null"`
	Desc       string `json:"desc"`
}

func (s *Storage) GetDatabaseById(id string) (*Db, bool, error) {
	database := &Db{}
	err := s.db.Where("id = ?", id).First(database).Error
	if gorm.IsRecordNotFoundError(err) {
		return database, false, nil
	}
	return database, true, err
}

func (s *Storage) UpdateDatabase(database *Db) error {
	return s.db.Save(database).Error
}

func (s *Storage) DelDatabase(database *Db) error {
	return s.db.Delete(database).Error
}

func (s *Storage) GetDatabases() ([]Db, error) {
	databases := []Db{}
	err := s.db.Find(&databases).Error
	return databases, err
}
