样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    KEY index_column_a (column_a),  -- 先创建索引
    PRIMARY KEY (id)  -- 然后创建主键约束
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```
