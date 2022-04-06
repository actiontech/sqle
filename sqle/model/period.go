package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
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

func (r *Periods) SelfCheck() bool {
	for _, p := range *r {
		periodStartTime, err := time.Parse(FormatHourAndMinute, fmt.Sprintf("%02d:%02d", p.StartHour, p.StartMinute))
		if err != nil {
			return false
		}
		periodStopTime, err := time.Parse(FormatHourAndMinute, fmt.Sprintf("%02d:%02d", p.EndHour, p.EndMinute))
		if err != nil {
			return false
		}
		if periodStopTime.Before(periodStartTime) || periodStopTime.Equal(periodStartTime) {
			return false
		}
	}
	return true
}

const FormatHourAndMinute = "15:04"

func (r *Periods) IsWithinScope(executeTime time.Time) bool {
	et, err := time.Parse(FormatHourAndMinute, executeTime.Format(FormatHourAndMinute))
	if err != nil {
		return false
	}
	for _, period := range *r {
		periodStartTime, err := time.Parse(FormatHourAndMinute, fmt.Sprintf("%02d:%02d", period.StartHour, period.StartMinute))
		if err != nil {
			continue
		}
		periodStopTime, err := time.Parse(FormatHourAndMinute, fmt.Sprintf("%02d:%02d", period.EndHour, period.EndMinute))
		if err != nil {
			continue
		}
		if et.After(periodStopTime) || et.Before(periodStartTime) {
			continue
		}
		return true
	}
	return false
}
