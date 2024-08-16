//go:build !enterprise
// +build !enterprise

package server

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/model"
)

func newPushJob(p *model.ReportPushConfig) (PushJob, error) {
	return nil, fmt.Errorf("type %v not supported", p.Type)
}
