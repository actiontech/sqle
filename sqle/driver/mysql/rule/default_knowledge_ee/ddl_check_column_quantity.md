样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_1 INT,...,column_39 INT,  -- column总数不超过阈值，默认值：40
    PRIMARY KEY (id),
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
