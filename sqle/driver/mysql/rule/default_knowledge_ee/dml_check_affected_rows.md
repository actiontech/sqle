不建议使用：

```
UPDATE table_a SET column_a=1 WHERE ...  -- 先用SELECT count(*) FROM table_a WHERE ... 查看行数
```
```
DELETE FROM table_a WHERE ...  -- 先用SELECT count(*) FROM table_a WHERE ... 查看行数
```
