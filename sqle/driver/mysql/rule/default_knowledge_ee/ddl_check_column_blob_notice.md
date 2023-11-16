样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a BLOB DEFAULT NULL COMMENT 'column_a',  -- 不建议column使用BLOB类型
    column_b TEXT DEFAULT NULL COMMENT 'column_b',  -- 不建议column使用TEXT类型
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```