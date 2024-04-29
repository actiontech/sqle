//go:build !enterprise
// +build !enterprise

package optimization

func getDefaultRulesKnowledge() (map[string]string, error) {
	return nil, nil
}
