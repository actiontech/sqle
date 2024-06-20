子查询嵌套层数过多可能会导致性能下降和查询执行时间延长，当嵌套的子查询层数增加时，查询引擎需要逐层执行子查询，将子查询的结果作为父查询的条件，这会导致查询的复杂度呈指数增长，还会增加数据库的内存消耗和磁盘IO操作。

样例说明：

```
SELECT 
  t_a.id,(SELECT id FROM table_d WHERE id=1) AS d_id  -- 子查询作为列 
FROM
  table_a AS t_a 
JOIN
  (SELECT id FROM table_b WHERE id>=1) AS t_b ON t_a.id=t_b.id  -- 子查询作为表
WHERE 
  t_a.id IN (SELECT id FROM table_c WHERE id>=1)  -- 子查询作为表达式
-- 日常子查询作为表和表达式嵌套较多，嵌套层数不能超过阈值3层
```