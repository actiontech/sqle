使用整数类型（如INT、BIGINT、SMALLINT）的shardkey字段可以提供更高效的数据路由和查询计划生成，因为整数类型的比较和排序操作相对较快。此外，字符类型（如CHAR或VARCHAR）的shardkey字段可以提供更灵活的数据划分和查询条件的表达，以满足不同的业务需求。

```
/*! TDDL:SHARD_KEY(shard_id) */
CREATE TABLE table_a (
    shard_id INT, -- BIGINT , SMALLINT , CHAR , VARCHAR
    column_a INT,
    column_b VARCHAR(10)
) ENGINE=TDB
```

