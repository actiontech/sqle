//go:build !enterprise
// +build !enterprise

package rule

func GetDefaultRulesKnowledge() (map[string]string, error) {
	return nil, nil
}
