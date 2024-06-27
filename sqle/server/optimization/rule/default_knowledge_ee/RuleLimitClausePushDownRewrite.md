重写说明：

规则描述  
Limit子句下推优化通过尽可能的 “下压” Limit子句，提前过滤掉部分数据, 减少中间结果集的大小，减少后续计算需要处理的数据量, 以提高查询性能。

譬如如下的案例，在外查询有一个Limit子句，可以将其下推至内层查询执行：
```
select *
from (select c_nationkey nation, 'C' as type, count(1) num
      from customer
      group by c_nationkey
      union
      select s_nationkey nation, 'S' as type, count(1) num
      from supplier
      group by nation) as nation_s
order by nation limit 20, 10
```
重写之后的SQL如下:
```
select *
from (
(select customer.c_nationkey as nation, 'C' as ````type````, count(1) as num
        from customer
        group by customer.c_nationkey
        order by customer.c_nationkey limit 30)
       union
(select supplier.s_nationkey as nation, 'S' as ````type````, count(1) as num
  from supplier
  group by supplier.s_nationkey
  order by supplier.s_nationkey limit 30)) as nation_s
order by nation_s.nation limit 20, 10
```

触发条件  
外查询有一个LIMIT子句  
外查询没有GROUP BY子句  
外查询的FROM只有一个表引用，且是一个子查询  
外查询没有其他条件  
子查询为单个查询或是UNION/UNION ALL连接的多个子查询（或者是一个外连接的外表）  
OFFSET的值小于指定阈值
