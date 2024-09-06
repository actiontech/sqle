重写说明：

规则描述  
IN子查询是指符合下面形式的子查询，IN子查询可以改写成等价的相关EXISTS子查询或是内连接，从而可以产生一个新的过滤条件，如果该过滤条件上有合适的索引，或是通过索引推荐引擎推荐合适的索引，可以获得更好的性能。
```
(expr1, expr2...) [NOT] IN (SELECT expr3, expr4, ...)
```
IN子查询重写为EXISTS  
譬如下面的IN子查询语言是为了获取最近一年内有订单的用户信息，
```
select * from customer where c_custkey in (select o_custkey from orders where O_ORDERDATE>=current_date - interval 1 year)
```
它可以重写为exists子查询，从而可以产生一个过滤条件（c_custkey = o_custkey）：
```
select * from customer where exists (select * from orders where c_custkey = o_custkey and O_ORDERDATE>=current_date - interval 1 year)
```
IN子查询重写为内关联  
如果子查询的查询结果是不重复的，则IN子查询可以重写为两个表的关联，从而让数据库优化器可以规划更优的表连接顺序。

譬如下面的SQL， c_custkey是表customer的主键，
```
select * from orders where o_custkey in (select c_custkey from customer where c_phone like '139%')
```
则上面的查询语句可以重写为
```
select orders.* from orders, customer where o_custkey=c_custkey and c_phone like '139%'
```

触发条件  
如果子查询的结果集是不重复的，可以重写为内关联
