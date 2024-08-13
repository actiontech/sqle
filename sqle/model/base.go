// base.go file contains :
// - public function of Models;
// - public const variables.
package model

import (
	"database/sql/driver"
	"encoding/json"
)

const (

	// used by Model:
	// User, UserGroup
	Enabled  = 0
	Disabled = 1
)

type Strings []string

func (t *Strings) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Strings) Value() (driver.Value, error) {
	return json.Marshal(t)
}
