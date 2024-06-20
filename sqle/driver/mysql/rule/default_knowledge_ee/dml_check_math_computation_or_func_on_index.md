如果对索引列使用了数学运算或函数，会改变其原有的数据结构和排序方式，导致无法使用索引进行快速查询。

样例说明：

```
SELECT 
  column_a
FROM 
  table_a
WHERE
  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00' AND column_b+3>0  -- 例：column_a, column_b为索引列
```

重写说明：

规则描述  
在索引列上的运算将导致索引失效，容易造成全表扫描，产生严重的性能问题。所以需要尽量将索引列上的运算转换到常量端进行，譬如下面的SQL。
```
select * from tpch.orders where adddate(o_orderdate,  INTERVAL 31 DAY) =date '2019-10-10'
```
adddate函数将导致o_orderdate上的索引不可用，可以将其转换成下面这个等价的SQL，以便使用索引提升查询效率。
```
select * from tpch.orders where o_orderdate = subdate(date '2019-10-10' , INTERVAL 31 DAY);
```
可以帮助转换大量的函数以及+、-、*、/运算符相关的操作。点击获取该优化的更详细信息。

触发条件  
过滤条件是个AND过滤条件（非连接条件）  
过滤条件是个可索引条件  
在索引列上存在计算或函数  
