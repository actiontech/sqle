重写说明：

规则描述  
如果一个查询中既包含来自同一个表的排序字段也包含分组字段，但字段顺序不同，可以通过调整分组字段顺序，使其和排序字段顺序一致，这样数据库可以避免一次排序操作。

考虑以下两个SQL, 二者唯一的不同点是分组字段的顺序（第一个SQL是c_custkey, c_name, 第二个SQL是c_name,c_custkey），由于分组字段中不包括grouping set/cube/roll up等高级grouping操作，所以两个SQL是等价的。但是二者的执行计划及执行效率却不一样，因此可以考虑将第一个SQL重写为第二个SQL。
```
select o_custkey, o_orderdate, sum(O_TOTALPRICE)
from orders
group by o_custkey,o_orderdate
order by o_orderdate;
```
重写为：
```
select o_custkey, o_orderdate, sum(o_totalprice)
from orders
group by o_orderdate,o_custkey
order by o_orderdate;
```
触发条件  
在一个QueryBlock中存在成员大于1的order子句及group子句  
子句中引用的是同一个数据表中的列且无函数或计算  
order子句中的列是group子句的真子集  
order子句不是group子句的前缀
