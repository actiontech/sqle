根据最左前缀原则，应该将最常用的查询条件字段放在组合索引的最左侧位置，这样可以最大程度地利用索引的优势，提高查询效率。

样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',  -- 例：选择性第二
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',  -- 例：选择性第一
    column_c DATE DEFAULT NULL COMMENT 'column_c',  -- 例：选择性第三
    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',
    INDEX index_column_b_a_c (column_b, column_a, column_c), 
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
