样例说明：

```
SELECT column_a FROM table_a WHERE column_a IS NULL  --不建议使用NULL

SELECT column_a FROM table_a WHERE column_a IS NOT NULL --不建议使用NOT NULL
```