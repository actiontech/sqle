大字段的存储和查询会占用较多的资源，如果将其与其他字段存放在同一张表中，会导致整张表的性能下降。而将大字段单独存放在一张表中，可以减少对整张表的影响，提高查询效率。

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) NOT NULL DEFAULT '' COMMENT 'column_b',
    PRIMARY KEY (id),  
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
```
CREATE TABLE table_b (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    table_a_id INT NOT NULL DEFAULT 0 COMMENT '关联table_a的id字段',  -- 跟主表的主键做关联关系
    column_t TEXT COMMENT 'column_t',
    PRIMARY KEY (id),
    KEY index_a_id (table_a_id)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_b'
```
