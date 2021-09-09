// pmm-agent
// Copyright 2019 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agents

import (
	"context"

	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
)

// Change represents built-in Agent status change and/or QAN collect request.
type Change struct {
	Status        inventorypb.AgentStatus
	MetricsBucket []*agentpb.MetricsBucket
}

// BuiltinAgent is a common interface for all built-in Agents.
type BuiltinAgent interface {
	// Run extracts stats data and sends it to the channel until ctx is canceled.
	Run(ctx context.Context)

	// Changes returns channel that should be read until it is closed.
	Changes() <-chan Change
}
