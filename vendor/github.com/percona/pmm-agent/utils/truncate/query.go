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

// Package truncate privides strings truncation utilities.
package truncate

var maxQueryLength = 2048

// Query limits passed query string to 2048 Unicode runes, truncating it if necessary.
func Query(q string) (query string, truncated bool) {
	runes := []rune(q)
	if len(runes) <= maxQueryLength {
		return string(runes), false
	}

	// copy MySQL behavior
	return string(runes[:maxQueryLength-4]) + " ...", true
}
