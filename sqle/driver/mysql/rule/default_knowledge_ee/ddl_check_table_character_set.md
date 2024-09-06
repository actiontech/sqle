样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a'   --建议指定字符集
```