package inspector

import (
	"fmt"
	"sqle/model"
)

func (i *Inspect) Advise(rules []model.Rule) error {
	defer i.closeDbConn()
	i.Logger().Info("start advise sql")
	for _, commitSql := range i.Task.CommitSqls {
		currentSql := commitSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			for _, rule := range rules {
				i.currentRule = rule
				if handler, ok := RuleHandlerMap[rule.Name]; ok {
					if handler.Func == nil {
						continue
					}
					for _, node := range sql.Stmts {
						err := handler.Func(i, node)
						if err != nil {
							return err
						}
					}
				}
			}
			currentSql.InspectStatus = model.TASK_ACTION_DONE
			currentSql.InspectLevel = i.Results.level()
			currentSql.InspectResult = i.Results.message()

			// print osc
			oscCommandLine, err := i.generateOSCCommandLine(sql.Stmts[0])
			if err != nil {
				return err
			}
			if oscCommandLine != "" {
				if currentSql.InspectResult != "" {
					currentSql.InspectResult += "\n"
				}
				currentSql.InspectResult = fmt.Sprintf("%s[osc]%s",
					currentSql.InspectResult, oscCommandLine)
			}

			// clean up results
			i.Results = newInspectResults()
			return nil
		})
		if err != nil {
			i.Logger().Error("add commit sql to task failed")
			return err
		}
	}
	err := i.Do()
	if err != nil {
		i.Logger().Error("advise sql failed")
	} else {
		i.Logger().Info("advise sql finish")
	}
	return err
}
