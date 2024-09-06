重写说明：

规则描述  
滤条件下推（FPPD）通过尽可能的 “下压” 过滤条件至SQL中的内部查询块，提前过滤掉部分数据, 减少中间结果集的大小，进而减少后续计算需要处理的数据量，进而提升SQL执行性能，FPPD属于重写优化。

譬如如下的案例中，在外查询有一个条件nation = 100，可以下压到personDT子查询中：
```
select *
from (select c_nationkey nation, 'C' as type, count(1) num
      from customer
      group by nation
      union
      select s_nationkey nation, 'S', count(1) num
      from supplier
      group by nation) as person
where nation = 100
```
重写之后的SQL如下:
```
select *
from (select c_nationkey nation, 'C' as type, count(1) num
      from customer
      where c_nationkey = 100
      group by nation
      union
      select s_nationkey nation, 'S', count(1) num
      from supplier
      where s_nationkey = 100
      group by nation) as person
```

触发条件  
过滤条件是个AND过滤条件（非连接条件）  
过滤条件的字段来自FROM子查询（如果是视图，应该被视图定义的SQL替换掉）  
该子查询没有被 查询折叠优化消除掉  
该子查询本身没有LIMIT子句  
该子查询本身没有rownum或rank等窗口函数
