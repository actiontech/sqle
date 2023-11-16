当引用和操作变量时，查询优化器无法在编译时确定变量的具体值，这导致了查询计划的不确定性。不确定的查询计划可能会导致性能下降或查询结果的不准确。

样例说明：
禁止以下变量用法
```
/*! TDDL:SHARD_KEY(shard_id) */
SET @variable_a = 1
SET @variable_b = @variable_a+1
SELECT column_a,column_b FROM table_a WHERE column_a=@variable_b
UPDATE table_a SET column_b=1 WHERE column_a=@variable_b
```