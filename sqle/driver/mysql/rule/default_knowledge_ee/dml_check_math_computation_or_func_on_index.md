如果对索引列使用了数学运算或函数，会改变其原有的数据结构和排序方式，导致无法使用索引进行快速查询。

样例说明：
```
SELECT 
  column_a
FROM 
  table_a
WHERE
  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00' AND column_b+3>0  -- 例：column_a, column_b为索引列
```