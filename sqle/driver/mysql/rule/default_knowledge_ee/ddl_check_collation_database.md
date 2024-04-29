```
CREATE DATABASE db_a CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci

CREATE TABLE table_a (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',
    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',
    PRIMARY KEY (id),
    KEY index_a (column_a)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a' 
-- 如果使用当前库默认字符集和字符集排序，默认（DEFAULT后） 部分可以不指定
```
重写说明：

```
规则描述
ORDER BY 子句中的所有表达式需要按统一的 ASC 或 DESC 方向排序，才能利用索引来避免排序；如果ORDER BY 语句对多个不同条件使用不同方向的排序无法使用索引。

触发条件
有多个排序字段
存在两种排序方向
```