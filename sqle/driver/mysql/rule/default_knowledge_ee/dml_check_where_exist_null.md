样例说明：

```
SELECT column_a FROM table_a WHERE column_a IS NULL  --不建议使用NULL
```
```
SELECT column_a FROM table_a WHERE column_a IS NOT NULL --不建议使用NOT NULL
```

重写说明：

规则描述  
= null并不能判断表达式为空,= null总是被判断为假。判断表达式为空应该使用is null.  
case expr when nulll也并不能判断表达式为空, 判断表达式为空应该case when expr is null。在where/having的筛选条件的错误写法还比较容易发现并纠正，而在藏在case 语句里使用null值判断就比较难以被发现。

触发条件  
语句中存在 = null 或是case when expr is null判断逻辑
