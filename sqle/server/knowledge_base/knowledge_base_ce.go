//go:build !enterprise
// +build !enterprise

package knowledge_base

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/model"
)

func LoadKnowledge(rulesMap map[string][]*model.Rule) error {
	return nil
}

func CheckKnowledgeBaseLicense() error {
	return fmt.Errorf("license not support knowledge base")
}
