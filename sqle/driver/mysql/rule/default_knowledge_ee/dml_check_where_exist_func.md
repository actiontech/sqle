样例说明：

```
SELECT 
  column_a
FROM 
  table_a
WHERE
  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00' --不建议column_a条件字段使用函数
```