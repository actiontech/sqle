package pganalyze_collector

import (
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NullTimeToNullTimestamp(in null.Time) *NullTimestamp {
	if !in.Valid {
		return &NullTimestamp{Valid: false}
	}

	ts := timestamppb.New(in.Time)

	return &NullTimestamp{Valid: true, Value: ts}
}
