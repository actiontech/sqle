package model

import (
	"fmt"
	"strconv"
	"sync"

	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"github.com/jinzhu/gorm"
)

type SqlWhitelist struct {
	Model
	Value         string `json:"value" gorm:"not null;type:text"`
	Desc          string `json:"desc"`
	MessageDigest string `json:"message_digest" gorm:"type:char(32) not null comment 'md5 data';" `
}

func (s SqlWhitelist) TableName() string {
	return "sql_whitelist"
}
func (s *Storage) GetSqlWhitelistItemById(sqlWhiteId string) (*SqlWhitelist, bool, error) {
	sqlWhitelist := &SqlWhitelist{}
	err := s.db.Table("sql_whitelist").Where("id = ?", sqlWhiteId).First(sqlWhitelist).Error
	if err == gorm.ErrRecordNotFound {
		return sqlWhitelist, false, nil
	}
	return sqlWhitelist, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
func (s *Storage) GetSqlWhitelist(pageIndex, pageSize int) ([]SqlWhitelist, uint32, error) {
	var count uint32
	sqlWhitelist := []SqlWhitelist{}
	if pageSize == 0 {
		err := s.db.Order("id desc").Find(&sqlWhitelist).Count(&count).Error
		return sqlWhitelist, count, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	err := s.db.Model(&SqlWhitelist{}).Count(&count).Error
	if err != nil {
		return sqlWhitelist, 0, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	err = s.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&sqlWhitelist).Error
	return sqlWhitelist, count, errors.New(errors.CONNECT_STORAGE_ERROR, err)

}

var sqlWhitelistMD5Map map[string] /*whitelist id*/ string /*whitelist message digest md5Data*/

var sqlWhitelistMutex sync.Mutex

func (s *SqlWhitelist) InitSqlWhitelistMD5Map() error {
	sqlWhitelistMD5Map = make(map[string]string, 0)
	storage := GetStorage()
	sqlWhitelist, _, err := storage.GetSqlWhitelist(0, 0)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("init sql whitelist error %v", err))
		return err
	}
	sqlWhitelistMutex.Lock()
	for _, v := range sqlWhitelist {
		sqlWhitelistMD5Map[strconv.Itoa(int(v.ID))] = v.MessageDigest
	}
	sqlWhitelistMutex.Unlock()
	return nil
}
func (s *SqlWhitelist) PutSqlWhitelistMD5() {
	if len(sqlWhitelistMD5Map) == 0 {
		go s.InitSqlWhitelistMD5Map()
		return
	}
	sqlWhitelistMutex.Lock()
	sqlWhitelistMD5Map[strconv.Itoa(int(s.ID))] = s.MessageDigest
	sqlWhitelistMutex.Unlock()
}

func (s *SqlWhitelist) RemoveSqlWhitelistMD5() {
	if len(sqlWhitelistMD5Map) == 0 {
		go s.InitSqlWhitelistMD5Map()
		return
	}
	sqlWhitelistMutex.Lock()
	delete(sqlWhitelistMD5Map, strconv.Itoa(int(s.ID)))
	sqlWhitelistMutex.Unlock()
}

func GetSqlWhitelistMD5Map() map[string]struct{} {
	if len(sqlWhitelistMD5Map) == 0 {
		if err := (&SqlWhitelist{}).InitSqlWhitelistMD5Map(); err == nil {
			return getSqlWhitelistMD5Map()
		}
		return nil
	}
	return getSqlWhitelistMD5Map()
}
func getSqlWhitelistMD5Map() map[string]struct{} {
	ret := make(map[string]struct{})
	sqlWhitelistMutex.Lock()
	for _, md5Data := range sqlWhitelistMD5Map {
		ret[md5Data] = struct{}{}
	}
	sqlWhitelistMutex.Unlock()
	return ret
}
