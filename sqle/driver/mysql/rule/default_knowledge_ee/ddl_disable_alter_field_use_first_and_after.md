不建议操作：

```
ALTER TABLE table_a ADD column_c INT NOT NULL DEFAULT 0 COMMENT 'column_c'  **FIRST**
```
```
ALTER TABLE table_a ADD column_d INT NOT NULL DEFAULT 0 COMMENT 'column_d' **AFTER id**
```