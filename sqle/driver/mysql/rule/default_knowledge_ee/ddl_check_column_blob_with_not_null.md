样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a BLOB NOT NULL COMMENT 'column_a',  -- 写入数据时又未对该字段指定值会导致写入失败
    column_b TEXT NOT NULL COMMENT 'column_b',  -- 写入数据时又未对该字段指定值会导致写入失败
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```