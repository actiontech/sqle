当违反规则时，result赋值的格式参考如下：
rulepkg.AddResult(input.Res, input.Rule, SQLE00025, util.JoinColumnNames(violateColumns))
rulepkg.AddResult(input.Res, input.Rule, SQLE00004)