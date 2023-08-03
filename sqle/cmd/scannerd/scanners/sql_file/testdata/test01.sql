-- SQL文件样例
-- 创建表
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    age INT
);

CREATE TABLE orders (
    order_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    product_name VARCHAR(100),
    amount DECIMAL(10, 2)
);

-- 插入数据
INSERT INTO users (name, age) VALUES ('Alice', 30);
INSERT INTO users (name, age) VALUES ('Bob', 25);
INSERT INTO users (name, age) VALUES ('Carol', 28);

INSERT INTO orders (user_id, product_name, amount) VALUES (1, 'Product A', 100.50);
INSERT INTO orders (user_id, product_name, amount) VALUES (1, 'Product B', 150.25);
INSERT INTO orders (user_id, product_name, amount) VALUES (2, 'Product C', 75.80);

-- 更新数据
UPDATE users SET age = 31 WHERE name = 'Alice';

-- 删除数据
DELETE FROM orders WHERE product_name = 'Product B';