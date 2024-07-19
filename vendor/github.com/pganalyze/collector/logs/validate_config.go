package logs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pganalyze/collector/state"
)

const MinSupportedLogMinDurationStatement = 10

func ValidateLogCollectionConfig(server *state.Server, settings []state.PostgresSetting) (bool, bool, bool, string) {
	var disabled = false
	var ignoreLogStatement = false
	var ignoreLogDuration = false
	var reasons []string

	if server.Config.DisableLogs {
		disabled = true
		reasons = append(reasons, "the collector setting disable_logs or environment variable PGA_DISABLE_LOGS is set")
	}

	if !disabled {
		for _, setting := range settings {
			if setting.Name == "log_min_duration_statement" && setting.CurrentValue.Valid {
				numVal, err := strconv.Atoi(setting.CurrentValue.String)
				if err != nil {
					continue
				}
				if numVal != -1 && numVal < MinSupportedLogMinDurationStatement {
					ignoreLogDuration = true
					reasons = append(reasons,
						fmt.Sprintf("log_min_duration_statement is set to '%d', below minimum supported threshold '%d'", numVal, MinSupportedLogMinDurationStatement),
					)
				}
			} else if setting.Name == "log_duration" && setting.CurrentValue.Valid {
				if setting.CurrentValue.String == "on" {
					ignoreLogDuration = true
					reasons = append(reasons, "log_duration is set to unsupported value 'on'")
				}
			} else if setting.Name == "log_statement" && setting.CurrentValue.Valid {
				if setting.CurrentValue.String == "all" {
					ignoreLogStatement = true
					reasons = append(reasons, "log_statement is set to unsupported value 'all'")
				}
			} else if setting.Name == "log_error_verbosity" && setting.CurrentValue.Valid {
				if setting.CurrentValue.String == "verbose" {
					disabled = true
					reasons = append(reasons, "log_error_verbosity is set to unsupported value 'verbose'")
				}
			}
		}
	}

	return disabled, ignoreLogStatement, ignoreLogDuration, strings.Join(reasons, "; ")
}
