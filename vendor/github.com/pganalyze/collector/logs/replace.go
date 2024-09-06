package logs

import (
	"sort"

	"github.com/pganalyze/collector/state"
)

type logRange struct {
	start int64
	end   int64 // When working with this range, the character at the index of end is *excluded*
}

const replacementChar = 'X'

// ReplaceSecrets - Replaces the secrets of the specified kind with the replacement character in the text
func ReplaceSecrets(input []byte, logLines []state.LogLine, filterLogSecret []state.LogSecretKind) []byte {
	var goodRanges []logRange

	filterUnidentified := false
	for _, k := range filterLogSecret {
		if k == state.UnidentifiedLogSecret {
			filterUnidentified = true
		}
	}

	for _, logLine := range logLines {
		goodRanges = append(goodRanges, logRange{start: logLine.ByteStart, end: logLine.ByteContentStart})
		if logLine.ReviewedForSecrets {
			sort.Slice(logLine.SecretMarkers, func(i, j int) bool {
				return logLine.SecretMarkers[i].ByteStart < logLine.SecretMarkers[j].ByteEnd
			})

			// We're creating a good range when we find a filtered secret or the end (everything before is marked as good)
			nextIdxToEvaluate := logLine.ByteContentStart
			for _, m := range logLine.SecretMarkers {
				filter := false
				for _, k := range filterLogSecret {
					if m.Kind == k {
						filter = true
					}
				}
				if filter {
					firstFilteredIdx := logLine.ByteContentStart + int64(m.ByteStart)
					goodRanges = append(goodRanges, logRange{start: nextIdxToEvaluate, end: firstFilteredIdx})
					nextIdxToEvaluate = logLine.ByteContentStart + int64(m.ByteEnd)
				}
			}
			// No more markers means the rest of the line is safe
			if nextIdxToEvaluate < logLine.ByteEnd {
				goodRanges = append(goodRanges, logRange{start: nextIdxToEvaluate, end: logLine.ByteEnd})
			}
		} else if !filterUnidentified {
			goodRanges = append(goodRanges, logRange{start: logLine.ByteContentStart, end: logLine.ByteEnd})
		}
	}
	sort.Slice(goodRanges, func(i, j int) bool {
		return goodRanges[i].start < goodRanges[j].start
	})

	lastGood := int64(0)
	for _, r := range goodRanges {
		for i := lastGood; i < r.start; i++ {
			input[i] = replacementChar
		}
		lastGood = r.end
	}
	if len(goodRanges) > 0 || filterUnidentified {
		for i := lastGood; i < int64(len(input)); i++ {
			input[i] = replacementChar
		}
	}
	return input
}
