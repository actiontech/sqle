重写说明：

规则描述  
对于仅进行存在性测试的子查询,如果子查询包含DISTINCT通常可以删除,以避免一次去重操作，譬如

IN子查询:
```
SELECT * FROM customer WHERE c_custkey IN (SELECT DISTINCT o_custkey FROM orders);
```
可以简化为:
```
SELECT * FROM customer WHERE c_custkey IN (SELECT o_custkey FROM orders);
```

触发条件  
使用IN/EXISTS子查询进行存在性判断  
子查询中存在DISTINCT/DISTINCT/UNIQUE关键字
