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

package version

import (
	"regexp"
	"strings"
)

// regexps to extract version numbers from the `SELECT version()` output
var (
	postgresDBRegexp = regexp.MustCompile(`PostgreSQL ([\d\.]+)`)
)

func ParsePostgreSQLVersion(v string) string {
	m := postgresDBRegexp.FindStringSubmatch(v)
	if len(m) != 2 {
		return ""
	}

	parts := strings.Split(m[1], ".")
	switch len(parts) {
	case 1: // major only
		return parts[0]
	case 2: // major and patch
		return parts[0]
	case 3: // major, minor, and patch
		return parts[0] + "." + parts[1]
	default:
		return ""
	}
}
