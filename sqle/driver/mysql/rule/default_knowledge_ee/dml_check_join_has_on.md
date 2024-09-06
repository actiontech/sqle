样例说明：

```
SELECT 
  t_a.column_a,t_b.column_b
FROM 
  table_a AS t_a
JOIN
  table_b AS t_b ON t_a.column_a=t_b.column_b  -- 如果没有关联条件，那JOIN就没意义，只是查询2张无关联表。
```
