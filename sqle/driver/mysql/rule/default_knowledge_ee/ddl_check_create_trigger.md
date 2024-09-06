样例说明：

```
CREATE TRIGGER trigger_example
AFTER INSERT ON table_a
FOR EACH ROW
BEGIN -- 触发器内容，违反规则
    INSERT INTO table_b (column_a) VALUES ('xxx');
END"
```
