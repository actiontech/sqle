// base.go file contains :
// - public function of Models;
// - public const variables.
package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (

	// used by Model:
	// User, UserGroup
	Enabled  = 0
	Disabled = 1
)

type Strings []string

func (t *Strings) Scan(input interface{}) error {
	if input == nil {
		return nil
	}
	if data, ok := input.([]byte); !ok {
		return fmt.Errorf("strings Scan input is not bytes")
	} else {
		return json.Unmarshal(data, t)
	}
}

func (t Strings) Value() (driver.Value, error) {
	return json.Marshal(t)
}
