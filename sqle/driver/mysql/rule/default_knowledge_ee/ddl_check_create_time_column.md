使用CREATE_TIME字段，有利于问题查找跟踪和检索数据，同时避免后期对数据生命周期管理不便 ，默认值为CURRENT_TIMESTAMP可保证时间的准确性


样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT '' COMMENT 'column_b',
    CREATE_TIME TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间-版本控制',  -- CURRENT_TIMESTAMP当前服务器的日期时间
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```