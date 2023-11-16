样例说明：

```
CREATE TABLE table_a (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id', --建议使用
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```