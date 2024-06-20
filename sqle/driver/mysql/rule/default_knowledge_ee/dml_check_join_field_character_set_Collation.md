索引是按照特定的字符集和排序规则进行存储和排序的，如果关联字段的字符集和排序规则不一致，会导致无法使用索引进行快速查询。

样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a'  -- 建议定义表级别的字符集和排序规则
```
```
CREATE TABLE table_b (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_b'  -- 建议定义表级别的字符集和排序规则
```
```
SELECT 
  t_a.column_a,t_b.column_b
FROM 
  table_a AS t_a
JOIN
  table_b AS t_b ON t_a.column_a=t_b.column_a 
```
