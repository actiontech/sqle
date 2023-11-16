样例说明：

```
SELECT 
  column_a
FROM 
  table_a
WHERE
  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00'    AND column_b+3<>0  --不建议条件中使用函数
```
