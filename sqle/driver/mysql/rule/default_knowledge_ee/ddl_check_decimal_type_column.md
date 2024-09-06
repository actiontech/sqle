样例说明：
```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b DECIMAL(5,2) DEFAULT '0.00' COMMENT 'column_b',  -- DECIMAL(5,2)代表总位数5，整数部分5-2=3位，超出会报错，小数部分四舍五入保留2位
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```