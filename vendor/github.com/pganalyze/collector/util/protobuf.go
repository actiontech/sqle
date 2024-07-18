package util

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Timestamp timestamppb.Timestamp

func (ts *Timestamp) Scan(value interface{}) error {
	if ts != nil {
		return fmt.Errorf("Can't scan timestamp into nil reference")
	}

	var t time.Time
	var protoTs *timestamppb.Timestamp
	if value == nil {
		return nil
	}

	t = value.(time.Time)
	protoTs = timestamppb.New(t)

	*ts = Timestamp(*protoTs)

	return nil
}
