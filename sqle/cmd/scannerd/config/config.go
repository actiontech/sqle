package config

import "fmt"

type ScannerType int32

const (
	ScannerTypeSlowQuery ScannerType = 0
	ScannerTypeMyBatis   ScannerType = 1
)

func (t ScannerType) String() string {
	switch t {
	case ScannerTypeMyBatis:
		return "MyBatis"
	case ScannerTypeSlowQuery:
		return "SlowQuery"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

type Config struct {
	Host          string
	Port          string
	Dir           string
	AuditPlanName string
	Token         string
	Typ           ScannerType
}
