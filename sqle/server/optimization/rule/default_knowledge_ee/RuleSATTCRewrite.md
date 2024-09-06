重写说明：

规则描述  
SAT-TC(SATisfiability-Transitive Closure) 重写优化是指分析一组相关的查询条件，去发现是否有条件自相矛盾、简化或是推断出新的条件，从而帮助数据库优化器选择更好的执行计划，提升SQL性能。

考虑下面的例子，
```
select c.c_name FROM customer c where c.c_name = 'John' and c.c_name = 'Jessey'
```
由于条件自相矛盾，所以重写后的SQL为，
```
select c.c_name from customer as c where  1 = 0
```

触发条件  
谓词间存在矛盾(例如 c_custkey=1 AND c_custkey=0),或者  
可以从谓词集中推断出新的谓词(例如 c_custkey=1 AND c_custkey=o_custkey 意味着 o_custkey=1)。  
谓词可以简化(例如 c_custkey <> c_custkey or c_name = 'b' 可以简化为 c_name = 'b')
