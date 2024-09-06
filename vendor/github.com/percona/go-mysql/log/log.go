/*
Copyright (c) 2019, Percona LLC.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package log provides an interface and data structures for MySQL log parsers.
// Log parsing yields events that are aggregated to calculate metric statistics
// like max Query_time. See also percona.com/go-mysql/event/.
package log

import (
	"time"
)

// An event is a query like "SELECT col FROM t WHERE id = 1", some metrics like
// Query_time (slow log) or SUM_TIMER_WAIT (Performance Schema), and other
// metadata like default database, timestamp, etc. Metrics and metadata are not
// guaranteed to be defined--and frequently they are not--but at minimum an
// event is expected to define the query and Query_time metric. Other metrics
// and metadata vary according to MySQL version, distro, and configuration.
type Event struct {
	Offset        uint64    // byte offset in file at which event starts
	OffsetEnd     uint64    // byte offset in file at which event ends
	Ts            time.Time // timestamp of event
	Admin         bool      // true if Query is admin command
	Query         string    // SQL query or admin command
	User          string
	Host          string
	Db            string
	Server        string
	LabelsKey     []string
	LabelsValue   []string
	TimeMetrics   map[string]float64 // *_time and *_wait metrics
	NumberMetrics map[string]uint64  // most metrics
	BoolMetrics   map[string]bool    // yes/no metrics
	RateType      string             // Percona Server rate limit type
	RateLimit     uint               // Percona Server rate limit value
}

// NewEvent returns a new Event with initialized metric maps.
func NewEvent() *Event {
	event := new(Event)
	event.TimeMetrics = make(map[string]float64)
	event.NumberMetrics = make(map[string]uint64)
	event.BoolMetrics = make(map[string]bool)
	return event
}

// Options encapsulate common options for making a new LogParser.
type Options struct {
	StartOffset        uint64                                // byte offset in file at which to start parsing
	FilterAdminCommand map[string]bool                       // admin commands to ignore
	Debug              bool                                  // print trace info to STDERR with standard library logger
	Debugf             func(format string, v ...interface{}) // use this function for logging instead of log.Printf (Debug still should be true)
	DefaultLocation    *time.Location                        // DefaultLocation to assume for logs in MySQL < 5.7 format.
}

// A LogParser sends events to a channel.
type LogParser interface {
	Start() error
	Stop()
	EventChan() <-chan *Event
}
