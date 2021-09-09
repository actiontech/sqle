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

package perfschema

import (
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

func getSummaries(q *reform.Querier) (map[string]*eventsStatementsSummaryByDigest, error) {
	rows, err := q.SelectRows(eventsStatementsSummaryByDigestView, "WHERE DIGEST IS NOT NULL AND DIGEST_TEXT IS NOT NULL")
	if err != nil {
		return nil, errors.Wrap(err, "failed to query events_statements_summary_by_digest")
	}
	defer rows.Close() //nolint:errcheck

	res := make(map[string]*eventsStatementsSummaryByDigest)
	for {
		var ess eventsStatementsSummaryByDigest
		if err = q.NextRow(&ess, rows); err != nil {
			break
		}

		// From https://dev.mysql.com/doc/relnotes/mysql/8.0/en/news-8-0-11.html:
		// > The Performance Schema could produce DIGEST_TEXT values with a trailing space. […] (Bug #26908015)
		*ess.DigestText = strings.TrimSpace(*ess.DigestText)

		res[*ess.Digest] = &ess
	}
	if err != reform.ErrNoRows {
		return nil, errors.Wrap(err, "failed to fetch events_statements_summary_by_digest")
	}
	return res, nil
}

// summaryCache provides cached access to performance_schema.events_statements_summary_by_digest.
// It retains data longer than this table.
type summaryCache struct {
	retain time.Duration

	rw    sync.RWMutex
	items map[string]*eventsStatementsSummaryByDigest
	added map[string]time.Time
}

// newSummaryCache creates new summaryCache.
func newSummaryCache(retain time.Duration) *summaryCache {
	return &summaryCache{
		retain: retain,
		items:  make(map[string]*eventsStatementsSummaryByDigest),
		added:  make(map[string]time.Time),
	}
}

// get returns all current items.
func (c *summaryCache) get() map[string]*eventsStatementsSummaryByDigest {
	c.rw.RLock()
	defer c.rw.RUnlock()

	res := make(map[string]*eventsStatementsSummaryByDigest, len(c.items))
	for k, v := range c.items {
		res[k] = v
	}
	return res
}

// refresh removes expired items in cache, then adds current items.
func (c *summaryCache) refresh(current map[string]*eventsStatementsSummaryByDigest) {
	c.rw.Lock()
	defer c.rw.Unlock()

	now := time.Now()

	for k, t := range c.added {
		if now.Sub(t) > c.retain {
			delete(c.items, k)
			delete(c.added, k)
		}
	}

	for k, v := range current {
		c.items[k] = v
		c.added[k] = now
	}
}
