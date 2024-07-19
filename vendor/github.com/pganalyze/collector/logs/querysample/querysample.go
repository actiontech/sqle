package querysample

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/guregu/null"
	"github.com/pganalyze/collector/logs/util"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
)

func TransformAutoExplainToQuerySample(logLine state.LogLine, explainText string, queryRuntime string) (state.PostgresQuerySample, error) {
	queryRuntimeMs, _ := strconv.ParseFloat(queryRuntime, 64)
	if strings.HasPrefix(explainText, "{") { // json format
		if util.WasTruncated(explainText) {
			return state.PostgresQuerySample{}, fmt.Errorf("auto_explain output was truncated and can't be parsed as JSON")
		} else {
			return transformExplainJSONToQuerySample(logLine, explainText, queryRuntimeMs)
		}
	} else if strings.HasPrefix(explainText, "Query Text:") { // text format
		return transformExplainTextToQuerySample(logLine, explainText, queryRuntimeMs)
	} else {
		return state.PostgresQuerySample{}, fmt.Errorf("unsupported auto_explain format")
	}
}

func transformExplainJSONToQuerySample(logLine state.LogLine, explainText string, queryRuntimeMs float64) (state.PostgresQuerySample, error) {
	var explainJSONOutput state.ExplainPlanContainer

	if err := json.Unmarshal([]byte(explainText), &explainJSONOutput); err != nil {
		return state.PostgresQuerySample{}, err
	}

	// Remove query text from EXPLAIN itself, to avoid duplication and match EXPLAIN (FORMAT JSON)
	sampleQueryText := strings.TrimSpace(explainJSONOutput.QueryText)
	explainJSONOutput.QueryText = ""

	return state.PostgresQuerySample{
		Query:             sampleQueryText,
		RuntimeMs:         queryRuntimeMs,
		OccurredAt:        logLine.OccurredAt,
		Username:          logLine.Username,
		Database:          logLine.Database,
		LogLineUUID:       logLine.UUID,
		HasExplain:        true,
		ExplainSource:     pganalyze_collector.QuerySample_AUTO_EXPLAIN_EXPLAIN_SOURCE,
		ExplainFormat:     pganalyze_collector.QuerySample_JSON_EXPLAIN_FORMAT,
		ExplainOutputJSON: &explainJSONOutput,
	}, nil
}

var autoExplainTextPlanDetailsRegexp = regexp.MustCompile(`^Query Text: (.+)\s+([\s\S]+)`)

func transformExplainTextToQuerySample(logLine state.LogLine, explainText string, queryRuntimeMs float64) (state.PostgresQuerySample, error) {
	explainParts := autoExplainTextPlanDetailsRegexp.FindStringSubmatch(explainText)
	if len(explainParts) != 3 {
		return state.PostgresQuerySample{}, fmt.Errorf("auto_explain output doesn't match expected format")
	}
	return state.PostgresQuerySample{
		Query:             strings.TrimSpace(explainParts[1]),
		RuntimeMs:         queryRuntimeMs,
		OccurredAt:        logLine.OccurredAt,
		Username:          logLine.Username,
		Database:          logLine.Database,
		LogLineUUID:       logLine.UUID,
		HasExplain:        true,
		ExplainSource:     pganalyze_collector.QuerySample_AUTO_EXPLAIN_EXPLAIN_SOURCE,
		ExplainFormat:     pganalyze_collector.QuerySample_TEXT_EXPLAIN_FORMAT,
		ExplainOutputText: explainParts[2],
	}, nil
}

func TransformLogMinDurationStatementToQuerySample(logLine state.LogLine, queryText string, queryRuntime string, queryProtocolStep string, parameterParts [][]string) (s state.PostgresQuerySample, ok bool) {
	// Ignore bind/parse steps of extended query protocol, since they are not the actual execution
	// See https://www.postgresql.org/docs/current/protocol-flow.html#PROTOCOL-FLOW-EXT-QUERY
	if queryProtocolStep == "bind" || queryProtocolStep == "parse" {
		return state.PostgresQuerySample{}, false
	}

	queryText = strings.TrimSpace(queryText)
	if queryText == "" {
		return state.PostgresQuerySample{}, false
	}

	sample := state.PostgresQuerySample{
		Query:       queryText,
		OccurredAt:  logLine.OccurredAt,
		Username:    logLine.Username,
		Database:    logLine.Database,
		LogLineUUID: logLine.UUID,
	}
	sample.RuntimeMs, _ = strconv.ParseFloat(queryRuntime, 64)
	for _, part := range parameterParts {
		if len(part) == 3 {
			if part[1] == "NULL" {
				sample.Parameters = append(sample.Parameters, null.NewString("", false))
			} else {
				sample.Parameters = append(sample.Parameters, null.StringFrom(part[2]))
			}
		}
	}
	return sample, true
}
