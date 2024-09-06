重写说明：

规则描述  
在索引列上的运算将导致索引失效，容易造成全表扫描，产生严重的性能问题。所以需要尽量将索引列上的运算转换到常量端进行，譬如下面的SQL。
```
select * from tpch.orders where adddate(o_orderdate,  INTERVAL 31 DAY) =date '2019-10-10'
```
`adddate`函数将导致`o_orderdate`上的索引不可用，可以将其转换成下面这个等价的SQL，以便使用索引提升查询效率。
```
select * from tpch.orders where o_orderdate = subdate(date '2019-10-10' , INTERVAL 31 DAY);
```

触发条件  
过滤条件是个`AND`过滤条件（非连接条件）  
过滤条件是个可索引条件  
在索引列上存在计算或函数