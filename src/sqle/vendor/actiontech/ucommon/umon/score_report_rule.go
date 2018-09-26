package umon

import (
	"regexp"
)

type ScoreReportRule struct {
	ExactlyMatches     map[string]string  `json:"exactly_matches"`
	RegexMatches       map[string]string  `json:"regex_matches"`
	Cron               string             `json:"cron"`
	PreservationNumber int                `json:"preservation_number"`
	Children           []*ScoreReportRule `json:"children"`
	CreatedAt          int64              `json:"created_at"`
}

/*
	return
		1. first child if match
		2. this if no child match
		3. null if no match
*/
func (u *ScoreReportRule) Match(tags map[string]string) *ScoreReportRule {
	if nil != u.ExactlyMatches {
		for tag, val := range u.ExactlyMatches {
			if "" != tags[tag] && val != tags[tag] {
				return nil
			}
		}
	}
	if nil != u.RegexMatches {
		for tag, val := range u.RegexMatches {
			re, err := regexp.Compile(val)
			if nil != err {
				return nil
			}
			if "" != tags[tag] && re.MatchString(tags[tag]) {
				return nil
			}
		}
	}
	if nil != u.Children {
		for _, child := range u.Children {
			if m := child.Match(tags); nil != m {
				return m
			}
		}
	}
	return u
}
