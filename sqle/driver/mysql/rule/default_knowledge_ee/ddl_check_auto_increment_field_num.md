样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',  -- AUTO_INCREMENT字段只能设置一个
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```