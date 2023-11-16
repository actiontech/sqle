样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a CHAR(30) DEFAULT NULL COMMENT 'column_a',  -- 不建议使用CHAR类型
    column_b VARCHAR(30) DEFAULT NULL COMMENT 'column_b',-- 建议使用VARCHAR类型
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
