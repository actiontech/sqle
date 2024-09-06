当使用DELETE/UPDATE...LIMIT语句时，只对部分数据进行删除或更新操作。这可能导致分片表在不同节点上的数据不一致，因为只有部分节点上的数据受到了影响，而其他节点上的数据保持不变。这会破坏数据的一致性，导致查询结果不准确。

样例说明：

```
DELETE FROM table_a LIMIT 10；  --分片表不建议使用DELETE....LIMIT
UPDATE table_a SET column_a=1 LIMIT 10；  --分片表不建议使用UPDATE....LIMIT
```
