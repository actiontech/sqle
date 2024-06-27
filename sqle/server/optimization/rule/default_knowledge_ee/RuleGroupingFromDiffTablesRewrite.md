重写说明：

规则描述  
如果分组字段来自不同的表，数据库优化器将没有办法利用索引的有序性来避免一次排序。如果WHERE或是HAVING子句里存在等值条件，可以排序或分组字段进行替换，使其来自同一张表，从而能够利用索引来避免一次排序。譬如下面的查询
```
select o_custkey, c_name, sum(o.O_TOTALPRICE) from customer c, orders o where o_custkey = c_custkey group by o_custkey, c_name
```
分组字段o_custkey, c_name来自两个表，且存在过滤条件o_custkey = c_custkey，可以重写为
```
select c_custkey, c_name, sum(o.O_TOTALPRICE) from customer c, orders o where o_custkey = c_custkey  group by c_custkey, c_name
```

触发条件  
GROUPBY字段来自不同表  
过滤条件是个可索引条件  
在索引列上不存在计算或函数
