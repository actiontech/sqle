禁止对分片键使用函数是为了确保数据的路由是可预测和确定的。函数的使用可能导致分片键的值被改变，从而无法准确地路由数据到正确的节点。

样例说明：

```
禁止以下分片键使用函数方式
/*! TDDL:SHARD_KEY(shard_id) */
SELECT column_a,column_b FROM table_a WHERE DATE(shard_id)='2020-01-01' AND ... 
INSERT INTO table_a(DATE(shard_id),column_a,column_b) VALUES (1,1,'a1') 
DELETE FROM table_a WHERE DATE(shard_id)='2020-01-01' AND ...  
UPDATE table_a SET column_a=1 WHERE DATE(shard_id)='2020-01-01' AND ...  
```