样例说明：

```
SELECT 
  column_a 
FROM 
  table_a 
WHERE 
  column_a=0 
GROUP BY column_b
HAVING column_b >1  --不建议条件放HAVING 中
```

重写说明：

规则描述  
从逻辑上，HAVING条件是在分组之后执行的，而WHERE子句上的条件可以在表访问的时候（索引访问）,或是表访问之后、分组之前执行，这两种条件都比在分组之后执行代价要小。

考虑下面的例子，
```
select c_custkey, count(*) from customer group by c_custkey having c_custkey < 100
```
重写后的SQL为，
```
select c_custkey, count(*) from customer where c_custkey < 100 group by c_custkey
```

触发条件  
HAVING子句中不存在聚集函数
