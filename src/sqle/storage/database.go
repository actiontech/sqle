package storage

import "github.com/jinzhu/gorm"

const (
	DB_TYPE_MYSQL = iota
	DB_TYPE_MYCAT
	DB_TYPE_SQLSERVER
)

type Db struct {
	gorm.Model
	DbType   int
	Alias    string
	Host     string
	Port     string
	User     string
	Password string
}

func (s *Storage) GetDatabaseById(id string) (*Db, error) {
	database := &Db{}
	err := s.db.First(database).Error
	return database, err
}

func (s *Storage) UpdateDatabase(database *Db) error {
	return s.db.Save(database).Error
}

func (s *Storage) DelDatabase(database *Db) error {
	return s.db.Delete(database).Error
}

func (s *Storage) GetDatabases() ([]*Db, error) {
	databases := []*Db{}
	err := s.db.Find(databases).Error
	return databases, err
}
