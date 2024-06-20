样例说明：

```
SELECT 
  t_a.column_a,t_b.column_b
FROM 
  table_a AS t_a
JOIN
  table_b AS t_b ON t_a.column_a=t_b.column_b  -- 例：column_a INT、column_b INT，建议匹配的字段类型一致，避免隐式转换
```
