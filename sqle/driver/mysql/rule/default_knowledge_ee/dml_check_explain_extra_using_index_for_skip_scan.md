不建议操作：

```
SELECT 
  column_a,column_b,column_c
FROM 
  table_a
ORDER BY
  column_a,column_b,column_c -- 当MySQL无法使用索引来满足ORDER BY子句的排序要求时，会使用文件排序
```
