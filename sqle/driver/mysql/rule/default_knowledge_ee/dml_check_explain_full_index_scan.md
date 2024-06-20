通过EXPLAIN查看是否使用了全索引扫描：

type字段：该字段表示访问表的方式。如果type的值是ALL，则表示MySQL将执行全表扫描，而不是使用索引进行查询。  
key字段：该字段显示MySQL选择的索引。如果key的值为NULL，则表示查询将执行全索引扫描。  

`如果type是ALL且key是NULL，则很可能发生了全索引扫描。`
