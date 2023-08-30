package machines

// DFATrans represents a Deterministic Finite Automatons state transition table
type DFATrans [][256]int

// DFAAccepting represents maps from accepting DFA states to match identifiers.
// These both identify which states are accepting states and which matches they
// belong to from the AST.
type DFAAccepting map[int]int

type lineCol struct {
	line, col int
}

// Compute the line and column of a particular index inside of a byte slice.
func mapLineCols(text []byte) []lineCol {
	m := make([]lineCol, len(text))
	line := 1
	col := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			col = 0
			line++
		} else {
			col++
		}
		m[i] = lineCol{line: line, col: col}
	}
	return m
}

// DFALexerEngine does the actual tokenization of the byte slice text using the
// DFA state machine. If the lexing process fails the Scanner will return
// an UnconsumedInput error.
func DFALexerEngine(startState, errorState int, trans DFATrans, accepting DFAAccepting, text []byte) Scanner {
	lineCols := mapLineCols(text)
	done := false
	matchID := -1
	matchTC := -1

	var scan Scanner
	scan = func(tc int) (int, *Match, error, Scanner) {
		if done && tc == len(text) {
			return tc, nil, nil, nil
		}
		startTC := tc
		if tc < matchTC {
			// we back-tracked so reset the last matchTC
			matchTC = -1
		} else if tc == matchTC {
			// the caller did not reset the tc, we are where we left
		} else if matchTC != -1 && tc > matchTC {
			// we skipped text
			matchTC = tc
		}
		state := startState
		for ; tc < len(text) && state != errorState; tc++ {
			if match, has := accepting[state]; has {
				matchID = match
				matchTC = tc
			}
			state = trans[state][text[tc]]
			if state == errorState && matchID > -1 {
				startLC := lineCols[startTC]
				endLC := lineCols[matchTC-1]
				match := &Match{
					PC:          matchID,
					TC:          startTC,
					StartLine:   startLC.line,
					StartColumn: startLC.col,
					EndLine:     endLC.line,
					EndColumn:   endLC.col,
					Bytes:       text[startTC:matchTC],
				}
				if matchTC == startTC {
					err := &EmptyMatchError{
						MatchID: matchID,
						TC:      tc,
						Line:    startLC.line,
						Column:  startLC.col,
					}
					return startTC, nil, err, scan
				}
				matchID = -1
				return matchTC, match, nil, scan
			}
		}
		if match, has := accepting[state]; has {
			matchID = match
			matchTC = tc
		}
		if startTC <= len(text) && matchID > -1 && matchTC == startTC {
			var startLC lineCol
			if startTC < len(text) {
				startLC = lineCols[startTC]
			}
			err := &EmptyMatchError{
				MatchID: matchID,
				TC:      tc,
				Line:    startLC.line,
				Column:  startLC.col,
			}
			matchID = -1
			return startTC, nil, err, scan
		}
		if startTC < len(text) && matchTC <= len(text) && matchID > -1 {
			startLC := lineCols[startTC]
			endLC := lineCols[matchTC-1]
			match := &Match{
				PC:          matchID,
				TC:          startTC,
				StartLine:   startLC.line,
				StartColumn: startLC.col,
				EndLine:     endLC.line,
				EndColumn:   endLC.col,
				Bytes:       text[startTC:matchTC],
			}
			matchID = -1
			return matchTC, match, nil, scan
		}
		if matchTC != len(text) && startTC >= len(text) {
			// the user has moved us farther than the text. Assume that was
			// the intent and return EOF.
			return tc, nil, nil, nil
		} else if matchTC != len(text) {
			done = true
			if matchTC == -1 {
				matchTC = 0
			}
			startLC := lineCols[startTC]
			etc := tc
			var endLC lineCol
			if etc >= len(lineCols) {
				endLC = lineCols[len(lineCols)-1]
			} else {
				endLC = lineCols[etc]
			}
			err := &UnconsumedInput{
				StartTC:     startTC,
				FailTC:      etc,
				StartLine:   startLC.line,
				StartColumn: startLC.col,
				FailLine:    endLC.line,
				FailColumn:  endLC.col,
				Text:        text,
			}
			return tc, nil, err, scan
		} else {
			return tc, nil, nil, nil
		}
	}
	return scan
}
