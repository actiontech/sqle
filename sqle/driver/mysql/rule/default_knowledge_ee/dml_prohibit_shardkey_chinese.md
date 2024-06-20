为了保持数据的一致性，shardkey的类型和内容应该是固定的。中文字符具有可变长度和复杂的编码方式，这可能导致数据划分和路由的不一致性。为了避免数据划分错误和查询结果的不准确，限制shardkey字段内容不包含中文。

```
/*! TDDL:SHARD_KEY(shard_id) */
不建议操作
INSERT INTO table_a(shard_id,column_a,column_b) VALUES ('张三',1,'a1')
```