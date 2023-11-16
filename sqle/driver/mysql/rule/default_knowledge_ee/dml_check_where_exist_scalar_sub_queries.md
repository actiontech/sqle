样例说明：

```
SELECT 
  column_a,(SELECT MAX(column_a) FROM table_b) AS max_value --不建议子查询中使用标量
FROM 
  table_a
WHERE
  column_a=1
```