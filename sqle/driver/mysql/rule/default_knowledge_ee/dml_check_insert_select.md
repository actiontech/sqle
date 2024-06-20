INSERT ... SELECT语句可能导致性能问题，特别是在大型数据集上，会导致数据库负载增加和执行时间延长，还可能导致不一致的数据插入从而破坏数据的一致性。

不建议操作：

```
INSERT INTO table_a(column_a,column_b,column_c) SELECT column_a,column_b,column_c FROM table_b
```