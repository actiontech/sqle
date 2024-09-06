//go:build enterprise
// +build enterprise

package rule

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

//go:embed default_knowledge_ee
var f embed.FS

const defaultKnowledgeRootDir = "default_knowledge_ee"

func getDefaultRulesKnowledge() (map[string]string, error) {
	res := make(map[string]string, 0)
	dir, err := f.ReadDir(defaultKnowledgeRootDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		ruleName := strings.TrimSuffix(entry.Name(), ".md")
		filePath := path.Join(defaultKnowledgeRootDir, entry.Name())
		content, err := f.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("read file [%v] failed: %v", filePath, err)
		}
		res[ruleName] = string(content)
	}
	return res, nil
}
