package params

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Periods []*Period

type Period struct {
	StartHour   int `json:"start_hour"`
	StartMinute int `json:"start_minute"`
	EndHour     int `json:"end_hour"`
	EndMinute   int `json:"end_minute"`
}

// Scan impl sql.Scanner interface
func (r *Periods) Scan(value interface{}) error {
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
	result := Periods{}
	err := json.Unmarshal(bytes, &result)
	*r = result
	return err
}

// Value impl sql.driver.Valuer interface
func (r Periods) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	v, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json value: %v", v)
	}
	return v, err
}

func (r *Periods) Copy() Periods {
	ps := make(Periods, 0, len(*r))
	for _, p := range *r {
		ps = append(ps, &Period{
			StartHour:   p.StartHour,
			StartMinute: p.StartMinute,
			EndHour:     p.EndHour,
			EndMinute:   p.EndMinute,
		})
	}
	return ps
}
