IN 、OR操作会导致查询无法走全索引，可将SQL拆分为多次等值查询。

样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    column_c DATE DEFAULT NULL COMMENT 'column_c',
    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',
    INDEX index_column_a_b_c (column_a, column_b, column_c),  --建立联合索引
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
```
SELECT 
  column_a,column_b,column_c,column_d 
FROM
  table_a 
WHERE 
  column_a IN (1,2,3) OR column_b='a'  --不建议使用IN 、OR条件进行查询
```
