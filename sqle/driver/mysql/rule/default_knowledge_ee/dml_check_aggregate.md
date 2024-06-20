使用聚合函数可能会导致性能问题，特别是在处理大量数据时，会引起不必要的计算开销，影响数据库的查询性能。

不建议使用：

```
SELECT COUNT(*) FROM table_a
```
```
SELECT SUM(column_a) FROM table_a
```
```
SELECT AVG(column_a) FROM table_a
```
```
SELECT MIN(column_a) FROM table_a
```
```
SELECT MAX(column_a) FROM table_a
```
```
SELECT GROUP_CONCAT(column_a) FROM table_a
```
```
SELECT column_a FROM table_a GROUP BY column_a HAVING column_b=0
```