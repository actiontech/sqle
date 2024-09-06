重写说明：

规则描述  
外连接优化指的是满足一定条件（外表具有NULL拒绝条件）的外连接可以转化为内连接，从而可以让数据库优化器可以选择更优的执行计划，提升SQL查询的性能。

考虑下面的例子，
```
select c_custkey from orders left join customer on c_custkey=o_custkey where C_NATIONKEY  < 20
```
C_NATIONKEY  < 20是一个customer表上的NULL拒绝条件，所以上面的左外连接可以重写为内连接，
```
select c_custkey from orders inner join customer on c_custkey=o_custkey where C_NATIONKEY  < 20
```

触发条件  
对于SQL，
```
SELECT * T1 FROM T1 LEFT JOIN T2 ON P1(T1,T2) WHERE P(T1,T2) AND R(T2)
```
如果R(T2) 是一个空拒绝条件条件(NFC)，那么以上的外连接可以转化为内连接，即
```
SELECT * T1 FROM T1 JOIN T2 ON P1(T1,T2) WHERE P(T1,T2) AND R(T2)
```
这样，优化器可以先应用R(T2) ，获取非常小的结果集，然后再和T1进行关联。
