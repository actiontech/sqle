重写说明：

规则描述  
如果有使用OR条件的UPDATE或DELETE语句，数据库优化器有可能无法使用索引来完成操作。
```
delete from lineitem where l_shipdate = date '2010-12-01' or l_partkey<100
```
如果这两个字段上都有索引，可以把它重写为多个DELETE语句，利用索引提升查询性能。
```
delete from lineitem where l_shipdate = date '2010-12-01';
```
```
delete from lineitem where l_partkey<100;
```

触发条件  
SQL为UPDATE或DELETE语句  
UPDATE或DELETE语句存在OR条件  
OR条件的各个分支都可以索引
