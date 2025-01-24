当违反规则时，result赋值的格式参考如下，注意SQLE00xxx前后没有双引号：
```go
rulepkg.AddResult(input.Res, input.Rule, SQLE00025, util.JoinColumnNames(violateColumns))
rulepkg.AddResult(input.Res, input.Rule, SQLE00004)
```

当获取规则参数时，格式参考如下：
```go
// 获取数值类型的规则参数
param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
if param == nil {
  return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
}

threshold := param.Int()
if threshold == 0 {
  return fmt.Errorf("param value should be greater than 0")
}

// 获取字符串类型的规则参数
param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
if param == nil {
  return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
}

requiredPrefix := param.String()
if requiredPrefix == "" {
  return fmt.Errorf("param value should not be empty")
}
```