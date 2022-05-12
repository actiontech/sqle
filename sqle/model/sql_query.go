package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type SqlQueryConfig struct {
	MaxPreQueryRows    int `json:"max_pre_query_rows"`
	QueryTimeoutSecond int `json:"query_timeout_second"`
}

// Scan impl sql.Scanner interface
func (c *SqlQueryConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal json value: %v", value)
	}
	if len(bytes) == 0 {
		return nil
	}
	result := SqlQueryConfig{}
	err := json.Unmarshal(bytes, &result)
	*c = result
	return err
}

// Value impl sql.driver.Valuer interface
func (c SqlQueryConfig) Value() (driver.Value, error) {
	v, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json value: %v", v)
	}
	return v, err
}

type SqlQueryHistory struct {
	Model
	RawSql string `json:"raw_sql" gorm:"type:text;not null"`
}

type SqlQueryExecutionSql struct {
	Model
	Sql        string `json:"sql" gorm:"type:text;not null"`
	InstanceId uint   `json:"instance_id" gorm:"not null"`
}
