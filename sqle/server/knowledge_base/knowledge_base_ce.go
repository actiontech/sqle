//go:build !enterprise
// +build !enterprise

package knowledge_base

import "github.com/actiontech/sqle/sqle/model"

func MigrateKnowledgeFromRules(rulesMap map[string][]*model.Rule) error{
	return nil
}
