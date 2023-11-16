为了将数据正确地路由到相应的节点，系统需要根据插入语句中的分片键值来确定数据应该插入到哪个分片节点上。如果插入字段不包含分片键，系统无法确定数据应该被路由到哪个节点，从而无法完成插入操作。

样例说明：

```
/*! TDDL:SHARD_KEY(shard_id) */
INSERT INTO table_a(shard_id,column_a,column_b) VALUES (1,1,'a1')
```