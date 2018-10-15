package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func NewMysql(user, password, host, port, schema string) (*Storage, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, schema))
	if err != nil {
		return nil, err
	}
	//db.LogMode(true)

	if err := createTable(db, &User{}); err != nil {
		return nil, err
	}
	if err := createTable(db, &Db{}); err != nil {
		return nil, err
	}
	if err := createTable(db, &Sql{}); err != nil {
		return nil, err
	}
	if err := createTable(db, &Task{}); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func createTable(db *gorm.DB, model interface{}) error {
	hasTable := db.HasTable(model)
	if db.Error != nil {
		return db.Error
	}
	if !hasTable {
		return db.CreateTable(model).Error
	}
	return nil
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) Exist(model interface{}) (bool, error) {
	var count int
	err := s.db.Model(model).Where(model).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Storage) Create(model interface{}) error {
	return s.db.Create(model).Error
}

func (s *Storage) Save(model interface{}) error {
	return s.db.Save(model).Error
}

func (s *Storage) Update(model interface{}, values map[string]interface{}) error {
	return s.db.Model(model).Update(values).Error
}
