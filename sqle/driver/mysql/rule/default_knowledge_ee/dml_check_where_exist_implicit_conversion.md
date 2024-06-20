样例说明：

```
SELECT 
  column_a
FROM 
  table_a
WHERE
  column_a='a'  -- 例：column_a是INT类型
```

重写说明：

规则描述  
当条件表达式的数据类型不同时，在查询执行过程中会进行一些隐式的数据类型转换。类型转换有时会应用于条件中的常量，有时会应用于条件中的列。当在列上应用类型转换时，在查询执行期间无法使用索引，可能导致严重的性能问题。譬如对于以下的SQL，
```
select count(*) from ORDERS where O_ORDERDATE = current_date();
```
如果O_ORDERDATE列的数据类型是CHAR(16)，那么O_ORDERDATE上的索引将不会被使用，导致全表扫描。解决方案通常有两个，一是ALTER TABLE改变O_ORDERDATE的数据类型，二是把current_date强制换换为CHAR类型。
```
 select count(*) ORDERS where ORDERS.O_ORDERDATE = cast(current_date() as CHAR(16));
```
触发条件  
条件表达式是个过滤条件，且是个可索引的过滤条件  
过滤条件两边的数据类型不一样  
根据数据库类型转换的优先级，数据库会优先转换列而非常量
