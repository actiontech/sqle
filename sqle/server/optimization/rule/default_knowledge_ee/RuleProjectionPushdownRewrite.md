重写说明：

规则描述  
投影下推指的通过删除DT子查询中无意义的列（在外查询中没有使用），来减少IO和网络的代价，同时提升优化器在进行表访问的规划时，采用无需回表的优化选项的几率。

考虑下面的例子，
```
SELECT count(1) FROM (SELECT c_custkey, avg(age) FROM customer group by c_custkey) AS derived_t1;
```
重写后的SQL为，
```
SELECT count(1) FROM (SELECT 1 FROM customer group by c_custkey) AS derived_t1;
```

触发条件  
内层选择列表中存在外层查询块没有使用的列
