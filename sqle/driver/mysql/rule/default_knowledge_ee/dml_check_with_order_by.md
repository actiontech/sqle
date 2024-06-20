样例说明：

```
UPDATE table_a SET column_a=1 WHERE column_a=0 ORDER BY column_b --UPDATE 语句不建议带ORDER BY
```
```
DELETE FROM table_a WHERE column_a=0 ORDER BY column_b --DELETE 语句不建议带ORDER BY
```
