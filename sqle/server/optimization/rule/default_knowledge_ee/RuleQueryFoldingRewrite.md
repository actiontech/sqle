重写说明：

规则描述  
查询折叠指的是把视图、CTE或是DT子查询展开，并与引用它的查询语句合并，来减少序列化中间结果集，或是触发更优的关于表连接规划的优化技术。

考虑下面的例子，
```
SELECT * FROM (SELECT c_custkey, c_name FROM customer) AS derived_t1;
```
重写后的SQL为，
```
SELECT c_custkey, c_name FROM customer
```

触发条件  
优化引擎针对不同的SQL语法结构，支持两种查询折叠的优化策略。其中第一种查询折叠的优化，MySQL 5.7以及PostgreSQL 14.0以上的版本都在优化器内部支持了此类优化；而第二类查询折叠的优化，在最新的MySQL及PostgreSQL版本中都没有支持。

查询折叠类型 I  
在视图本身中,没有distinct关键字  
在视图中没有聚集函数或窗口函数  
在视图本身中,没有LIMIT子句  
在视图本身中,没有UNION或者UNION ALL  
在外部查询块中,被折叠的视图不是外连接的一部分。

查询折叠类型 II  
在外部查询块中,视图是唯一的表引用  
在外部查询块中,没有分组、聚集函数和窗口函数  
在视图内部没有使用窗口函数
