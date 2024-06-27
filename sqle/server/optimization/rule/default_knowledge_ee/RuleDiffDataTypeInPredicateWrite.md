重写说明：

规则描述  
当条件表达式的数据类型不同时，在查询执行过程中会进行一些隐式的数据类型转换。类型转换有时会应用于条件中的常量，有时会应用于条件中的列。当在列上应用类型转换时，在查询执行期间无法使用索引，可能导致严重的性能问题。譬如对于以下的SQL，
```
select count(*) from ORDERS where O_ORDERDATE = current_date();
```
如果`O_ORDERDATE`列的数据类型是`CHAR(16)`，那么`O_ORDERDATE`上的索引将不会被使用，导致全表扫描。解决方案通常有两个，一是`ALTER TABLE`改变`O_ORDERDATE`的数据类型，二是把`current_date`强制换换为`CHAR`类型。
```
select count(*) ORDERS where ORDERS.O_ORDERDATE = cast(current_date() as CHAR(16));
```

触发条件  
条件表达式是个过滤条件，且是个可索引的过滤条件  
过滤条件两边的数据类型不一样  
根据数据库类型转换的优先级，数据库会优先转换列而非常量
