重写说明：

规则描述  
如果子查询没有LIMIT子句，那么子查询的排序操作就没有意义，可以将其删除而不影响最终的结果。一些案例如下：

EXISTS子查询
```
select * from lineitem l where exists (select * from part p where p.p_partkey=l.l_partkey and p.p_name = 'a' order by p_name )
```
可以重写为
```
select * from lineitem l where exists (select * from part p where p.p_partkey=l.l_partkey and p.p_name = 'a')
```

触发条件  
子查询存在ORDER子句  
子查询中没有LIMIT子句
