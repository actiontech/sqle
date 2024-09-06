样例说明：

```
SELECT 
  t_a.column_a
FROM 
  table_a AS t_a
JOIN 
  table_b AS t_b1 ON t_a.column_a=t_b1.column_a
JOIN --不建议单表多次连接
  table_b AS t_b2 ON t_a.column_a=t_b2.column_b
```