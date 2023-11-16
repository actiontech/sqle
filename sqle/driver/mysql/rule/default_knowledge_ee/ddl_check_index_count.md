样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    column_c DATE DEFAULT NULL COMMENT 'column_c',
    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',
    column_e TEXT DEFAULT NULL COMMENT 'column_e',
    INDEX index_column_a (column_a),  
    INDEX index_column_b (column_b),  
    INDEX index_column_c (column_c),  
    INDEX index_column_d (column_d), 
    INDEX index_column_e (column_e),  
    INDEX index_column_a_b (column_a, column_b),  -- 第六个索引，不建议超过阈值
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```