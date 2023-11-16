样例说明：

```
SELECT 
  t_a.column_a
FROM 
  table_a AS t_a
JOIN  
  table_b AS t_b ON t_a.column_a=t_b.column_a
JOIN
  table_c AS t_c ON t_a.column_a=t_c.column_a
JOIN  -- JOIN数量不超过阈值
  table_d AS t_d ON t_a.column_a=t_d.column_a

```
