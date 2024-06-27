重写说明：

规则描述
1. `= null`并不能判断表达式为空,`= null`总是被判断为假。判断表达式为空应该使用`is null`.

2. `case expr when nulll`也并不能判断表达式为空, 判断表达式为空应该`case when expr is null`。在where/having的筛选条件的错误写法还比较容易发现并纠正，而在藏在case 语句里使用null值判断就比较难以被发现。

触发条件  
语句中存在 `= null` 或是`case when expr is null`判断逻辑
