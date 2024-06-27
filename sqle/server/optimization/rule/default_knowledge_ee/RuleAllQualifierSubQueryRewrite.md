重写说明：

规则描述  
假设通过下面的SQL来获取订单系统关闭后注册的用户
```
select * from customer where c_regdate > all(select o_orderdate from orders)
```
如果子查询的结果中存在NULL，这个SQL永远返回为空。正确的写法应该是在子查询里加上非空限制，或使用max/min的写法
```
select * from customer where c_regdate > (select max(o_custkey) from orders)
```
推荐采用第二种写法，可以通过max/min重写进一步优化SQL。

触发条件  
ALL修饰的子查询条件
