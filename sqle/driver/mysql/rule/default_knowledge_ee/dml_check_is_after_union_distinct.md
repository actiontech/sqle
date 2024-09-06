样例说明：

```
SELECT column_a FROM table_a 
UNION ALL
SELECT column_a FROM table_b 
```

但要注意UNION ALL和UNION执行的结果是不一样的，UNION会去除重复数据，UNION ALL不会去除重复数据
