样例说明：

```
SELECT 
  column_b 
FROM 
  table_a 
WHERE 
  column_a=0 
GROUP BY column_b  
ORDER BY column_b  --GROUP BY语句中建议使用
```