样例说明：

```
SELECT 
  column_a
FROM
  table_a 
WHERE 
  column_a IN (SELECT column_a FROM table_b WHERE column_a>=1)  --不建议使用
```