重写说明：

规则描述  
ANY/SOME/ALL修饰的子查询来源自SQL-92 标准, 通常用于检查某个值与子查询返回的全部值或任意值的大小关系。使用ANY/SOME/ALL修饰的子查询执行效率低下,因为需要对子查询的结果集逐行进行比较,随着结果集大小增加而线性下降。可以通过查询重写的方式提升其执行效率。

譬如对于下面的SQL：
```
select * from orders where o_orderdate < all(select o_orderdate from orders where o_custkey > 100)
```
对于MySQL，可以重写为
```
select * from orders where o_orderdate < (select o_orderdate from orders where o_custkey > 100 order by o_orderdate asc limit 1)
```
对于PostgreSQL或Oracle，则可以重写为
```
select * from orders where o_orderdate < (select o_orderdate from orders where o_custkey > 100 order by o_orderdate asc nulls first limit 1)
```

触发条件  
SQL中存在ANY/SOME/ALL修饰的子查询
