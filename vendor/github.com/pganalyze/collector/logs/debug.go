package logs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
	uuid "github.com/satori/go.uuid"
)

func PrintDebugInfo(logFileContents string, logLines []state.LogLine, samples []state.PostgresQuerySample) {
	fmt.Printf("log lines: %d, query samples: %d\n", len(logLines), len(samples))
	groups := map[pganalyze_collector.LogLineInformation_LogClassification]int{}
	unclassifiedLogLines := []state.LogLine{}
	for _, logLine := range logLines {
		if logLine.ParentUUID != uuid.Nil {
			continue
		}

		groups[logLine.Classification]++

		if logLine.Classification == pganalyze_collector.LogLineInformation_UNKNOWN_LOG_CLASSIFICATION {
			unclassifiedLogLines = append(unclassifiedLogLines, logLine)
		}
	}

	for classification, count := range groups {
		fmt.Printf("%d x %s (%d)\n", count, classification, classification)
	}

	if len(unclassifiedLogLines) > 0 {
		fmt.Printf("\nUnclassified log lines:\n")
		for _, logLine := range unclassifiedLogLines {
			fmt.Printf("%s\n", logFileContents[logLine.ByteStart:logLine.ByteEnd])
			fmt.Printf("  Level: %s\n", logLine.LogLevel)
			fmt.Printf("  Content: %#v\n", logFileContents[logLine.ByteContentStart:logLine.ByteEnd])
			fmt.Printf("---\n")
		}
	}
}

func PrintDebugLogLines(logFileContents string, logLines []state.LogLine, classifications map[pganalyze_collector.LogLineInformation_LogClassification]bool) {
	fmt.Println("\nParsed log lines:")
	linesById := make(map[uuid.UUID]*state.LogLine)
	for _, logLine := range logLines {
		linesById[logLine.UUID] = &logLine
		if len(classifications) > 0 {
			var classifiedLine *state.LogLine
			if logLine.ParentUUID == uuid.Nil {
				classifiedLine = &logLine
			} else {
				classifiedLine = linesById[logLine.ParentUUID]
			}
			if _, ok := classifications[classifiedLine.Classification]; !ok {
				continue
			}
		}
		detailsStr, err := json.Marshal(logLine.Details)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", logFileContents[logLine.ByteStart:logLine.ByteEnd])
		fmt.Printf("  Level:          %s\n", logLine.LogLevel)
		if logLine.ParentUUID == uuid.Nil {
			fmt.Printf("  Classification: %s (%d)\n", logLine.Classification, logLine.Classification)
		}
		if len(logLine.Details) > 0 {
			fmt.Printf("  Details:        %s\n", detailsStr)
		}
		fmt.Printf("  Content:    %#v\n", logFileContents[logLine.ByteContentStart:logLine.ByteEnd])
		fmt.Printf("---\n")
	}
}

var HerokuPostgresDebugRegexp = regexp.MustCompile(`^(\w+ \d+ \d+:\d+:\d+ \w+ app\[postgres\] \w+ )?\[(\w+)\] \[\d+-\d+\] (.+)`)

type MaybeHerokuLogReader struct {
	LineReader
}

func NewMaybeHerokuLogReader(r io.Reader) *MaybeHerokuLogReader {
	return &MaybeHerokuLogReader{bufio.NewReader((r))}
}

func (lr *MaybeHerokuLogReader) ReadString(delim byte) (string, error) {
	line, err := lr.LineReader.ReadString(delim)
	if err != nil {
		return "", err
	}
	contentParts := HerokuPostgresDebugRegexp.FindStringSubmatch(line)
	if len(contentParts) == 4 {
		return contentParts[3], nil
	}

	return line, nil
}
