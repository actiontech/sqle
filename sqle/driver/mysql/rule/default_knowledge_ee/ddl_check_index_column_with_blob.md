样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b BLOB NOT NULL COMMENT 'column_b', 
    INDEX index_column_a (column_a),  
    INDEX index_column_b (column_b(50)),  -- 不建议BLOB类型的列索引
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';
```