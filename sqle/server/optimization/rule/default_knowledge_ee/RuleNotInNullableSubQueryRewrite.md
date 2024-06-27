重写说明：

规则描述  
对于以下想要查询没有订单用户的SQL，
```
select * from customer where c_custkey not in (select o_custkey from orders)
```
如果子查询的结果集里有空值，这个SQL永远返回为空。正确的写法应该是在子查询里加上非空限制，即
```
select * from customer where c_custkey not in (select o_custkey from orders where o_custkey is not null)
```
触发条件  
存在IN子查询条件  
IN子查询的选择列取值可以为NULL
