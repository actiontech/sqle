```
EXPLAIN SELECT column_a FROM table_a WHERE ...
-- 使用EXPLAIN查看扫描行数，超过10W的，筛选条件必须带上主键或者索引
```
