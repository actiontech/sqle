样例说明：

```
--不合并的方式，分开多次修改
ALTER TABLE table_a
ADD column_b INT NOT NULL DEFAULT 0 COMMENT 'column_b';
ALTER TABLE table_a
ADD INDEX index_column_a (column_a);
```
```
-- 建议使用以下合并的方式
ALTER TABLE table_a
ADD column_b INT NOT NULL DEFAULT 0 COMMENT 'column_b',
ADD INDEX index_column_a (column_a);"
```
