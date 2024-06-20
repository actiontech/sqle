样例说明：

```
SELECT 
  MAX(column_a)
FROM 
  table_a
WHERE
  DATEDIFF(column_b, column_c)  -- MySQL有很多内置函数，可根据需要调整
```
