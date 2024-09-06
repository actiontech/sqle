样例说明：

```
SELECT 
  column_b 
FROM 
  table_a 
WHERE 
  column_a=0 
GROUP BY column_b  
ORDER BY column_b  --GROUP BY语句中建议使用
```

重写说明：

规则描述  
在MySQL的早期版本中，即使没有order by子句，group by默认也会按分组字段排序，这就可能导致不必要的文件排序，影响SQL的查询性能。可以通过添加order by null来强制取消排序，禁用查询结果集的排序；

譬如下面的例子中
```
SELECT l_orderkey, sum(l_quantity)
FROM lineitem
GROUP BY l_orderkey;
```
在MySQL 5.x版本中，group by l_orderkey会引起默认排序, 可以通过添加order by null来避免该排序。
```
SELECT l_orderkey, sum(l_quantity)
FROM lineitem
GROUP BY l_orderkey
ORDER BY NULL;
```

触发条件  
MySQL数据库，版本低于8.0  
存在分组字段，且无排序字段
