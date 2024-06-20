样例说明：

```
SELECT 
  column_a,(SELECT MAX(column_a) FROM table_b) AS max_value --不建议子查询中使用标量
FROM 
  table_a
WHERE
  column_a=1
```

重写说明：

规则描述  
对于使用COUNT标量子查询来进行判断是否存在，可以重写为EXISTS子查询，从而避免一次聚集运算。譬如对于如下的SQL，
```
select * from customer where (select count(*) from orders where c_custkey=o_custkey) > 0
```
可以重写为,
```
select * from customer where exists(select 1 from orders where c_custkey=o_custkey)
```

触发条件  
存在COUNT标量子查询>0条件
