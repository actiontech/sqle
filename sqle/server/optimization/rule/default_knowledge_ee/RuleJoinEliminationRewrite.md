重写说明：

规则描述  
连接消除（Join Elimination）通过在不影响最终结果的情况下从查询中删除表，来简化SQL以提高查询性能。通常，当查询包含主键-外键连接并且查询中仅引用主表的主键列时，可以使用此优化。内连接和外连接都可以用于此重写优化。

内连接消除的案例
```
select o.* from orders o inner join customer c on c.c_custkey=o.o_custkey
```
订单表（orders）和客户表（customer）关联，且c_custkey是客户表的主键，那么客户表可以被消除掉，重写后的SQL如下：
```
select * from orders where o_custkey
```
外连接消除的案例
```
select o_custkey from orders left join customer on c_custkey=o_custkey
```
客户表可以被消除掉，重写后的SQL如下：
```
select orders.o_custkey from orders
```

触发条件  
查询包含主键-外键连接  
查询中仅引用主表的主键列
