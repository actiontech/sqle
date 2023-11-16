样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a BLOB COMMENT 'column_a',
    column_b TEXT COMMENT 'column_b',
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
```
INSERT INTO table_a(id) value(1)  -- 当插入数据不指定BLOB和TEXT类型字段时，字段值会被设置为NULL
```