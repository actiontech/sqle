/* DELIMITER str */
/* 注释中定义分隔符 预期不会被检测为分隔符语法 */
DELIMITER --good
-- DELIMITER str

/* 分隔符在注释中，预期不被分割 */
-- --good
/* --good */
DROP PROCEDURE IF EXISTS `_GS_GM_Check` --good
/* 分隔符在SQL语句中或注释中，预期不会被分割 */
/* 分隔符定义语法在SQL语句中，预期不会被识别 */
CREATE DELIMITER = `--good`@`%` PROCEDURE `DELIMITER`(vi_uid INT,vi_pwd VARCHAR(32),vi_ip VARCHAR(100),OUT vo_level INT,OUT vo_code INT)
BEGIN
/* 分隔符与语句粘连，预期会被分割 */
    END--good
/* 切换分隔符 */
DELIMITER ;
CREATE FUNCTION hello (s CHAR(20));
CREATE FUNCTION hello (s CHAR(20)) ;
-- 切换分隔符
DELIMITER //
use test //
/* 多条语句合并提交，预期切分结果会包含分隔符之间的多条sql */
SET @sql = 'SELECT * FROM employees WHERE salary = 2321.21';
SET @result = @sql;
PREPARE stmt FROM @result;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
CREATE DELIMITER = `//`@`%` PROCEDURE `DELIMITER`(vi_uid INT,vi_pwd VARCHAR(32),vi_ip VARCHAR(100),OUT vo_level INT,OUT vo_code INT)
BEGIN
    -- 分隔符与语句粘连，预期会被分割
    END;
//
DELIMITER ;

CREATE TABLE account (acct_num INT, amount DECIMAL(10,2)); /* ; */

CREATE TRIGGER ins_sum BEFORE /* ; */INSERT ON account
       FOR EACH ROW SET @sum = @sum + NEW.amount;

SET @sum = 0;

INSERT INTO account VALUES(137,14.98)/* ; */,(141,1937.50),(97,-100.00);/* ; */
SELECT @sum AS ';';/* ; */
DROP TRIGGER test.ins_sum;/* ; */
CREATE TRIGGER ins_transaction BEFORE INSERT ON account
    /* ; */  FOR EACH ROW PRECEDES ins_sum
       SET
       @deposits = @deposits + IF(NEW.amount>0,NEW.amount,0),
       @withdrawals = @withdrawals + IF(NEW.amount<0,-NEW.amount,0);
delimiter //
CREATE TRIGGER upd_check BEFORE UPDATE ON account
       FOR EACH ROW
       BEGIN
           IF NEW.amount < 0 THEN
               SET NEW.amount = 0;
           ELSEIF NEW.amount > 100 THEN
               SET NEW.amount = 100;
           END IF;
       END//
delimiter ;