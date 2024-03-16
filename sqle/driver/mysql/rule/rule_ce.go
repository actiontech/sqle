//go:build !enterprise
// +build !enterprise

package rule

func getDefaultRulesKnowledge() (map[string]string, error) {
	return nil, nil
}
