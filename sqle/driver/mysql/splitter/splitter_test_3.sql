-- 1
CREATE PROCEDURE nested_example()
BEGIN
    DECLARE v_outer_variable INT DEFAULT 0;

    -- 外层 BEGIN ... END 块
    BEGIN
        DECLARE v_inner_variable INT DEFAULT 10;

        -- 内层 BEGIN ... END 块
        BEGIN
            DECLARE v_nested_variable INT DEFAULT 20;
            
            IF v_outer_variable = 0 THEN
                SET v_outer_variable = v_inner_variable + v_nested_variable;
            END IF;
        END; -- 内层块结束
        WHILE v1 > 0 DO
            SET v1 = v1 - 1;
        END WHILE;
        IF v_outer_variable > 0 THEN
            SET v_outer_variable = v_outer_variable * 2;
        END IF;
    END; -- 外层块结束
END;
-- 2
CREATE PROCEDURE example_procedure()
BEGIN
    DECLARE v_variable INT DEFAULT 0;
    
    IF v_variable = 0 THEN
        SET v_variable = 1;
    ELSEIF v_variable = 1 THEN
        SET v_variable = 2;
    ELSE
        SET v_variable = 0;
    END IF;
    
    -- 其他处理语句
END;
-- 3
CREATE PROCEDURE doiterate(p1 INT)
BEGIN
  label1: LOOP
    SET p1 = p1 + 1;
    IF p1 < 10 THEN
      ITERATE label1;
    END IF;
    LEAVE label1;
  END LOOP label1;
  SET @x = p1;
END;
-- 4
CREATE PROCEDURE complex_loops_example()
BEGIN
    DECLARE v_counter INT DEFAULT 0;
    DECLARE v_limit INT DEFAULT 5;
    DECLARE v_sum INT DEFAULT 0;

    -- 使用 LOOP 结构
    loop_label: LOOP
        SET v_counter = v_counter + 1;
        
        -- 使用 WHILE 结构
        WHILE v_counter <= v_limit DO
            SET v_sum = v_sum + v_counter;
            
            -- 使用 REPEAT 结构
            REPEAT
                SET v_counter = v_counter + 1;
            UNTIL v_counter > v_limit
            END REPEAT;
        END WHILE;

        -- 退出 LOOP 的条件
        IF v_counter > v_limit THEN
            LEAVE loop_label;
        END IF;
    END LOOP;

    -- 输出结果
    SELECT v_sum AS total_sum;
END;

BEGIN;

delimiter %%
BEGIN%%