样例说明：

```
SELECT column_a FROM table_a WHERE column_a IN (NULL)  --不建议使用IN (NULL)

SELECT column_a FROM table_a WHERE column_a NOT IN (NULL) --不建议使用NOT IN (NULL)
```
