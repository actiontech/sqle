使用NOT NULL约束可以确保表中的每个记录都包含一个值，有助于维护数据的完整性，使数据的一致性更容易得到保证，数据库优化器也可以更好地执行查询优化，提高查询性能。

样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',  --建议使用NOT NULL约束
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',  --建议使用NOT NULL约束
    column_b VARCHAR(10) NOT NULL DEFAULT '' COMMENT 'column_b',--建议使用NOT NULL约束
    PRIMARY KEY (id),  
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```