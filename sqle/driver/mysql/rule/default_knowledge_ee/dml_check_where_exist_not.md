样例说明：

不建议使用以下查询条件
```
SELECT column_a FROM table_a WHERE column_a<>1
```
```
SELECT column_a FROM table_a WHERE column_a NOT IN (1,2)
```
```
SELECT column_a FROM table_a WHERE column_a NOT LIKE (1%)
```
```
SELECT column_a FROM table_a WHERE NOT EXISTS (SELECT column_a FROM table_b WHERE table_a.id=table_b.id)
```