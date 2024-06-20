样例说明：

大数据量的情况下，文件排序意味着SQL性能较低，会增加OS的开销，影响数据库性能
```
SELECT 
  column_a,column_b,column_c
FROM 
  table_a
ORDER BY
  column_a,column_b,column_c -- 当MySQL无法使用索引来满足ORDER BY子句的排序要求时，会使用文件排序
```