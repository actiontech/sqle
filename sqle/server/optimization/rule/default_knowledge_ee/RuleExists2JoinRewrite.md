重写说明：

规则描述  
EXISTS子查询向外层查询返回一个布尔值,表示是否存在满足条件的行。满足一定条件的EXISTS 子查询可以转换为 JOIN，从而可以让数据库优化器对驱动表的选择有更多的选择，从而生成更优的查询计划。

譬如对于如下的查询，
```
select * from lineitem l where exists (select * from part p where p.p_partkey=l.l_partkey and p.p_name = 'a')
```
如果子查询对于每一个l.l_partkey,都至多返回一行记录（即在等值条件的列上(p_partkey，p_name)有一个唯一性约束），则此子查询可以重写为如下的形式:
```
select l.* from lineitem as l, part as p where p.p_partkey = l.l_partkey and p.p_name = 'a'
```

触发条件  
EXISTS子查询条件由AND和其他条件关联  
EXISTS子查询无分组无LIMIT  
EXISTS子查询结果集返回UNIQUE的行  
EXISTS子查询和外查询关联方式为等值关联
