样例说明：

```
SELECT count(*) FROM table_a WHERE column_a=0  --建议使用 count(*)统计
```
```
SELECT count(column_a) FROM table_a WHERE column_a=0  --不建议使用 count(col)
```