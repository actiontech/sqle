样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b TIMESTAMP DEFAULT '1970-01-01 00:00:01' COMMENT 'column_b',  -- 1970-01-01 00:00:01是此字段类型最小时间，用UNIX_TIMESTAMP转换时为1
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```