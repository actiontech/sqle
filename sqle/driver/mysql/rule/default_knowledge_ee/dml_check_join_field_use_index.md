JOIN操作是基于索引的。如果JOIN字段没有索引，那么MySQL需要扫描整个表来找到匹配的行，这会导致查询性能下降。

样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    PRIMARY KEY (id),  
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
```
CREATE TABLE table_b (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_b INT NOT NULL DEFAULT 0 COMMENT 'column_b',
    PRIMARY KEY (id),  
    KEY index_b (column_b)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_b'
```
```
SELECT 
  t_a.column_a,t_b.column_b
FROM 
  table_a AS t_a
JOIN
  table_b AS t_b ON t_a.column_a=t_b.column_b  -- JOIN字段需要包含索引
```
