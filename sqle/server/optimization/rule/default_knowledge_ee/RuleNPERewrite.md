重写说明：


规则描述  
SQL的NPE(Null Pointer Exception)问题是指在SQL查询中,当聚合列全为NULL时,SUM、AVG等聚合函数会返回NULL,这可能会导致后续的程序出现空指针异常。譬如对于下面的SQL：
```
select sum(t.b) from (values row(1,null)) as t(a,b);
```
可以使用如下方式避免NPE问题:
```
SELECT IFNULL(SUM(t.b), 0) from (values row(1,null)) as t(a,b);
```
这会返回0而不是NULL,避免了空指针异常。

>   Oracle:NVL();  SQL Server和MS Access:ISNULL();  MySQL:IFNULL()或COALESCE();

触发条件  
1. SUM或AVG聚集函数  
2. 聚集函数的参数可能全为NULL, 包括  
   1. 参数是列，列定义可以为空  
   2. 参数是表达式，表达式可以为空  
   3. 列定义不可为空，但是是外连接的内表，结果可能为空
