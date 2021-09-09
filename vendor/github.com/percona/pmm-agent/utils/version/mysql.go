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

	"gopkg.in/reform.v1"
)

// regexps to extract version numbers from the `SHOW GLOBAL VARIABLES WHERE Variable_name = 'version'` output.
var (
	mysqlDBRegexp = regexp.MustCompile(`^\d+\.\d+`)
)

// GetMySQLVersion return parsed version of MySQL and vendor.
func GetMySQLVersion(q *reform.Querier) (string, string, error) {
	var name, ver string
	err := q.QueryRow(`SHOW /* pmm-agent:mysqlversion */ GLOBAL VARIABLES WHERE Variable_name = 'version'`).Scan(&name, &ver)
	if err != nil {
		return "", "", err
	}
	var ven string
	err = q.QueryRow(`SHOW /* pmm-agent:mysqlversion */ GLOBAL VARIABLES WHERE Variable_name = 'version_comment'`).Scan(&name, &ven)
	if err != nil {
		return "", "", err
	}

	version := mysqlDBRegexp.FindString(ver)
	var vendor string
	switch {
	case strings.Contains(strings.ToLower(name), "percona"):
		vendor = "percona"
	case strings.Contains(strings.ToLower(name), "mariadb"):
		vendor = "mariadb"
	default:
		vendor = "oracle"
	}

	return version, vendor, nil
}
