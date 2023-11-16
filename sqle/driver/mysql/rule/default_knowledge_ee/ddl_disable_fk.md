样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    FOREIGN KEY (column_a) REFERENCES other_table (other_column)   --不建议使用外检
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```