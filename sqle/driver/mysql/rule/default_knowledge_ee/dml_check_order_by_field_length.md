样例说明：

```
CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    long_text_column VARCHAR(2000) DEFAULT NULL COMMENT 'long_text_column',
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'
```
```
SELECT id FROM table_a ORDER BY long_text_column  -- 不建议使用，对长字段进行排序
```
