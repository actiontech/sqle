样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT DEFAULT 0 COMMENT 'column_a',  -- 没有非空约束
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    INDEX index_column_a (column_a),  --不建议索引字段没有非空约束
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```