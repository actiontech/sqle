重写说明：

规则描述  
对于使用MAX/MIN的子查询，
```
select * from customer where c_custkey = (select max(o_custkey) from orders)
```
可以重写为以下的形式，从而利用索引的有序来避免一次聚集运算，
```
select * from customer where c_custkey = (select o_custkey from orders order by o_custkey desc null last limit 1)
```

触发条件  
SQL中存在MAX/MIN的标量子查询
