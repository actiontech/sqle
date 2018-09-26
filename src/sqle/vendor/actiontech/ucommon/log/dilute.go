package log

import "time"

type dilutes map[string]*diluteRecord

type diluteRecord struct {
	firstTimestamp      time.Time
	lastTimestamp       time.Time
	checkpointTimestamp time.Time
}

var dilutesInstance = dilutes{}