重写说明：

规则描述  
如果使用OR条件的查询语句，数据库优化器有可能无法使用索引来完成查询。譬如，
```
select * from lineitem where l_shipdate = date '2010-12-01' or l_partkey<100
```
如果这两个字段上都有索引，可以把查询语句重写为UNION或UNION ALL查询，以便使用索引提升查询性能。
```
select * from lineitem where l_shipdate = date '2010-12-01'
union select * from lineitem where l_partkey<100
```
如果数据库支持INDEX MERGING（请参考如何创建高效的索引），也可以调整数据库相关参数启用INDEX MERGING优化策略来提升数据库性能。获取该优化的更详细信息。

触发条件  
OR连接的条件必须是可以利用索引的  
重写后的UNION语句估算代价比原SQL小  
如果OR分支的条件是互斥的，那么重写为UNION ALL
