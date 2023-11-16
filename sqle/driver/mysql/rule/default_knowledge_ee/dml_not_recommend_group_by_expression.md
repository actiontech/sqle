样例说明

```
SELECT 
  column_a 
FROM 
  table_a 
ORDER BY CASE WHEN column_b=3 THEN 1 ELSE column_b END --不建议使用表达式
```