```
CREATE DATABASE db_a CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
```
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a' 
-- 如果使用当前库默认字符集和字符集排序，默认（DEFAULT后） 部分可以不指定
```