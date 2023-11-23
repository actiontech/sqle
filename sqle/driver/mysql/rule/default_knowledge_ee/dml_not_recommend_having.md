样例说明：

```
SELECT 
  column_a 
FROM 
  table_a 
WHERE 
  column_a=0 
GROUP BY column_b
HAVING column_b >1  --不建议条件放HAVING 中
```