//go:build enterprise
// +build enterprise

package rule

import driverV2 "github.com/actiontech/sqle/sqle/driver/v2"

var defaultRuleKnowledgeMap = map[string]driverV2.RuleKnowledge{
	DDLCheckTableSize: {
		Content: "不建议操作：\n\n`ALTER TABLE -- 修改表结构`\n\n`DROP TABLE -- 删除表`\n\n`TRUNCATE TABLE -- 清空表数据`\n\n`RENAME TABLE -- 重命名表`\n",
	},
	DDLCheckIndexTooMany: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    column_c DATE DEFAULT NULL COMMENT 'column_c',\n    INDEX index_column_a (column_a),  \n    INDEX index_column_b_a (column_b, column_a),  \n    INDEX index_column_c_b_a (column_c, column_b, column_a),  -- a字段超过阈值\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	ConfigDMLExplainPreCheckEnable: {
		Content: "建议使用EXPLAIN+SQL语句验证\n\n样例说明：\n\n```\nEXPLAIN SELECT column_a FROM table_a WHERE column_a = 0;  \n```\n",
	},
	DDLCheckRedundantIndex: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a_b (column_a,column_b)，\n    KEY index_a (column_a)  --不建议使用\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	ConfigOptimizeIndexEnabled: {
		Content: "配置2个参数：\n\n```\n1.计算列基数阈值：在建立联合索引时，字段的顺序非常重要。将基数较高的列放到前面，可以走索引的过程中筛选掉更多的不需要的记录。\n    但是完全按照基数顺序建立联合索引，在一些场景下可能也不那么完美：\n    （1）表的记录较多，计算基数对性能影响较大\n    （2）表的数据分布可能经常发生变化，列的基数大小有可能因此发生变化\n    所以给出一个阈值，来控制联合索引是否按照基数从大到小的顺序建立，默认值是 1000000。\n2.联合索引最大列数：当联合索引的列数太多，可能会增加数据库对索引的维护成本。所以来控制联合索引的最大列数，默认值是 3。\n```\n",
	},
	DDLCheckPKWithoutIfNotExists: {
		Content: "样例说明：\n```\nCREATE TABLE IF NOT EXISTS table_a (  --建议使用 IF NOT EXISTS避免重复\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLCheckObjectNameLength: {
		Content: "样例说明：\n\n```\nCREATE TABLE this_is_a_64_character_table_name_000000000000000000000000000000 (\n    this_is_a_64_character_column_name_00000000000000000000000000000 INT,\n    INDEX this_is_a_64_character_index_name_000000000000000000000000000000 (this_is_a_64_character_column_name_00000000000000000000000000000)\n);-- 不能超过阈值\n```",
	},
	DDLCheckPKNotExist: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id)    --建表语句中需使用主键\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLCheckPKWithoutAutoIncrement: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',  --建议使用自增主键\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLCheckPKWithoutBigintUnsigned: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id', --建议使用\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DMLCheckJoinFieldType: {
		Content: "样例说明：\n```\nSELECT \n  t_a.column_a,t_b.column_b\nFROM \n  table_a AS t_a\nJOIN\n  table_b AS t_b ON t_a.column_a=t_b.column_b  -- 例：column_a INT、column_b INT，建议匹配的字段类型一致，避免隐式转换\n```\n",
	},
	DMLCheckJoinHasOn: {
		Content: "样例说明：\n```\nSELECT \n  t_a.column_a,t_b.column_b\nFROM \n  table_a AS t_a\nJOIN\n  table_b AS t_b ON t_a.column_a=t_b.column_b  -- 如果没有关联条件，那JOIN就没意义，只是查询2张无关联表。\n```\n",
	},
	DDLCheckColumnCharLength: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a CHAR(30) DEFAULT NULL COMMENT 'column_a',  -- 不建议使用CHAR类型\n    column_b VARCHAR(30) DEFAULT NULL COMMENT 'column_b',-- 建议使用VARCHAR类型\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DDLCheckFieldNotNUllMustContainDefaultValue: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n\n```\nINSERT INTO table_a(id) value(1) -- 如果不带DEFAULT值，INSERT时不包含该字段会报错\n```",
	},
	DDLDisableFK: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    FOREIGN KEY (column_a) REFERENCES other_table (other_column)   --不建议使用外检\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLDisableAlterFieldUseFirstAndAfter: {
		Content: "不建议操作：\n\nALTER TABLE table_a ADD column_c INT NOT NULL DEFAULT 0 COMMENT 'column_c'  **FIRST**\n\nALTER TABLE table_a ADD column_d INT NOT NULL DEFAULT 0 COMMENT 'column_d' **AFTER id**",
	},
	DDLCheckCreateTimeColumn: {
		Content: "使用CREATE_TIME字段，有利于问题查找跟踪和检索数据，同时避免后期对数据生命周期管理不便 ，默认值为CURRENT_TIMESTAMP可保证时间的准确性\n\n\n样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT '' COMMENT 'column_b',\n    CREATE_TIME TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间-版本控制',  -- CURRENT_TIMESTAMP当前服务器的日期时间\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckIndexCount: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    column_c DATE DEFAULT NULL COMMENT 'column_c',\n    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',\n    column_e TEXT DEFAULT NULL COMMENT 'column_e',\n    INDEX index_column_a (column_a),  \n    INDEX index_column_b (column_b),  \n    INDEX index_column_c (column_c),  \n    INDEX index_column_d (column_d), \n    INDEX index_column_e (column_e),  \n    INDEX index_column_a_b (column_a, column_b),  -- 第六个索引，不建议超过阈值\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLCheckUpdateTimeColumn: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT '' COMMENT 'column_b',\n    UPDATE_TIME TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间-版本控制',  -- 数据会随着UPDATE更新为当前服务器的日期时间\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DDLCheckCompositeIndexMax: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    column_c DATE DEFAULT NULL COMMENT 'column_c',\n    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',\n    INDEX index_column_a_b (column_a, column_b),  \n    INDEX index_column_b_a_c (column_b, column_a, column_c),  \n    INDEX index_column_c_b_a_d (column_c, column_b, column_a, column_d),  \n    INDEX index_column_d_b_c_a (column_d, column_b, column_c, column_a),  -- 第四个索引，不建议超过阈值\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```\n",
	},
	DDLCheckIndexNotNullConstraint: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT DEFAULT 0 COMMENT 'column_a',  -- 没有非空约束\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    INDEX index_column_a (column_a),  --不建议索引字段没有非空约束\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLCheckTableDBEngine: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'  --建议指定引擎，默认:INNODB\n```\n",
	},
	DDLCheckTableCharacterSet: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a'   --建议指定字符集\n```",
	},
	DDLCheckIndexedColumnWithBlob: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b BLOB NOT NULL COMMENT 'column_b', \n    INDEX index_column_a (column_a),  \n    INDEX index_column_b (column_b(50)),  -- 不建议BLOB类型的列索引\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DMLCheckWhereIsInvalid: {
		Content: "禁止使用：\n\n`SELECT column_a FROM table_a`\n\n`SELECT column_a FROM table_a WHERE 1`\n",
	},
	DDLCheckAlterTableNeedMerge: {
		Content: "样例说明：\n\n```\n--不合并的方式，分开多次修改\nALTER TABLE table_a\nADD column_b INT NOT NULL DEFAULT 0 COMMENT 'column_b';\nALTER TABLE table_a\nADD INDEX index_column_a (column_a);\n\n-- 建议使用以下合并的方式\nALTER TABLE table_a\nADD column_b INT NOT NULL DEFAULT 0 COMMENT 'column_b',\nADD INDEX index_column_a (column_a);\n```\n",
	},
	DMLDisableSelectAllColumn: {
		Content: "样例说明：\n\n```\nSELECT \n  *  --不建议使用通配符\nFROM \n  table_a\nWHERE\n  column_a=1\n```\n",
	},
	DDLDisableDropStatement: {
		Content: "样例说明：\n\n```\nDROP TABLE table_a;  --不建议使用，避免误删\n```\n",
	},
	DDLCheckTableWithoutComment: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a表注释' --建议添加表注释\n```\n",
	},
	DDLCheckColumnWithoutComment: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckColumnWithoutDefault: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT '' COMMENT 'column_b',  -- ''表示空字符串，NULL表示没有值，两者不同，查询时为where column_b=''\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckColumnTimestampWithoutDefault: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b TIMESTAMP DEFAULT '1970-01-01 00:00:01' COMMENT 'column_b',  -- 1970-01-01 00:00:01是此字段类型最小时间，用UNIX_TIMESTAMP转换时为1\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```",
	},
	DDLCheckColumnBlobWithNotNull: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a BLOB NOT NULL COMMENT 'column_a',  -- 写入数据时又未对该字段指定值会导致写入失败\n    column_b TEXT NOT NULL COMMENT 'column_b',  -- 写入数据时又未对该字段指定值会导致写入失败\n    PRIMARY KEY (id)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckColumnBlobDefaultIsNotNull: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a BLOB COMMENT 'column_a',\n    column_b TEXT COMMENT 'column_b',\n    PRIMARY KEY (id)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n```\nINSERT INTO table_a(id) value(1)  -- 当插入数据不指定BLOB和TEXT类型字段时，字段值会被设置为NULL\n```",
	},
	DDLCheckAutoIncrementFieldNum: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',  -- AUTO_INCREMENT字段只能设置一个\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckAllIndexNotNullConstraint: {
		Content: "```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),  -- 主键属于特殊索引\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DMLCheckSelectLimit: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a=1\nLIMIT 1000  -- 不超过阈值\n```\n",
	},
	DMLCheckWithOrderBy: {
		Content: "样例说明：\n\n```\nUPDATE table_a SET column_a=1 WHERE column_a=0 ORDER BY column_b --UPDATE 语句不建议带ORDER BY\n\nDELETE FROM table_a WHERE column_a=0 ORDER BY column_b --DELETE 语句不建议带ORDER BY\n```\n",
	},
	DMLCheckSelectWithOrderBy: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a=1\nORDER BY column_b  --SELECT 不建议使用ORDER BY\n```",
	},
	DMLCheckInsertColumnsExist: {
		Content: "样例说明：\n\n```\nINSERT INTO \n  table_a(column_a,column_b)  -- 不建议未明确指定列名\nVALUES\n  (1,'a1')\n```",
	},
	DMLCheckBatchInsertListsMax: {
		Content: "样例说明：\n\n```\nINSERT INTO \n  table_a(column_a,column_b) \nVALUES\n  (1,'a1'),\n  ...,\n  (100,'a100')  -- 不超过阈值\n```\n",
	},
	DMLCheckInQueryNumber: {
		Content: "样例说明：\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a IN (1,...,50)  -- 不超过阈值\n```\n",
	},
	DMLCheckWhereExistFunc: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00' --不建议column_a条件字段使用函数\n```",
	},
	DMLCheckWhereExistNot: {
		Content: "样例说明：\n\n```\n不建议使用以下查询条件\nSELECT column_a FROM table_a WHERE column_a<>1\nSELECT column_a FROM table_a WHERE column_a NOT IN (1,2)\nSELECT column_a FROM table_a WHERE column_a NOT LIKE (1%)\nSELECT column_a FROM table_a WHERE NOT EXISTS (SELECT column_a FROM table_b WHERE table_a.id=table_b.id)\n```",
	},
	DMLWhereExistNull: {
		Content: "样例说明：\n\n```\nSELECT column_a FROM table_a WHERE column_a IS NULL  --不建议使用NULL\n\nSELECT column_a FROM table_a WHERE column_a IS NOT NULL --不建议使用NOT NULL\n```",
	},
	DMLCheckWhereExistImplicitConversion: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a='a'  -- 例：column_a是INT类型\n```",
	},
	DMLCheckLimitMustExist: {
		Content: "样例说明：\n\nUPDATE table_a SET column_a=1 WHERE column_a=0 LIMIT 1000; -- 建议使用LIMIT",
	},
	DMLCheckWhereExistScalarSubquery: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a,(SELECT MAX(column_a) FROM table_b) AS max_value --不建议子查询中使用标量\nFROM \n  table_a\nWHERE\n  column_a=1\n```",
	},
	DDLCheckIndexesExistBeforeCreateConstraints: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    KEY index_column_a (column_a),  -- 先创建索引\n    PRIMARY KEY (id)  -- 然后创建主键约束\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a';\n```\n",
	},
	DMLCheckSelectForUpdate: {
		Content: "样例说明：\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a=1\nFOR UPDATE --不建议使用\n```",
	},
	DDLCheckDatabaseCollation: {
		Content: "```\nCREATE DATABASE db_a CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci\n\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a' \n-- 如果使用当前库默认字符集和字符集排序，默认（DEFAULT后） 部分可以不指定\n```",
	},
	DDLCheckDecimalTypeColumn: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b DECIMAL(5,2) DEFAULT '0.00' COMMENT 'column_b',  -- DECIMAL(5,2)代表总位数5，整数部分5-2=3位，超出会报错，小数部分四舍五入保留2位\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckBigintInsteadOfDecimal: {
		Content: "```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a BIGINT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DMLCheckSubQueryNestNum: {
		Content: "```\n子查询嵌套层数过多可能会导致性能下降和查询执行时间延长，当嵌套的子查询层数增加时，查询引擎需要逐层执行子查询，将子查询的结果作为父查询的条件，这会导致查询的复杂度呈指数增长，还会增加数据库的内存消耗和磁盘IO操作。\n```\n\n样例说明：\n\n```\nSELECT \n  t_a.id,(SELECT id FROM table_d WHERE id=1) AS d_id  -- 子查询作为列 \nFROM\n  table_a AS t_a \nJOIN\n  (SELECT id FROM table_b WHERE id>=1) AS t_b ON t_a.id=t_b.id  -- 子查询作为表\nWHERE \n  t_a.id IN (SELECT id FROM table_c WHERE id>=1)  -- 子查询作为表达式\n-- 日常子查询作为表和表达式嵌套较多，嵌套层数不能超过阈值3层\n```\n",
	},
	DMLCheckNeedlessFunc: {
		Content: "样例说明：\n```\nSELECT \n  MAX(column_a)\nFROM \n  table_a\nWHERE\n  DATEDIFF(column_b, column_c)  -- MySQL有很多内置函数，可根据需要调整\n```\n",
	},
	DMLCheckFuzzySearch: {
		Content: "**样例说明：**\n\n不建议操作\n\n`SELECT column_a FROM table_a WHERE column_a LIKE '%xxx%'`\n\n`SELECT column_a FROM table_a WHERE column_a LIKE '%xxx'`\n",
	},
	DMLCheckNumberOfJoinTables: {
		Content: "样例说明：\n\n```\nSELECT \n  t_a.column_a\nFROM \n  table_a AS t_a\nJOIN  \n  table_b AS t_b ON t_a.column_a=t_b.column_a\nJOIN\n  table_c AS t_c ON t_a.column_a=t_c.column_a\nJOIN  -- JOIN数量不超过阈值\n  table_d AS t_d ON t_a.column_a=t_d.column_a\n\n```\n",
	},
	DMLCheckIfAfterUnionDistinct: {
		Content: "样例说明：\n```\nSELECT column_a FROM table_a \nUNION ALL\nSELECT column_a FROM table_b \n```\n\n但要注意UNION ALL和UNION执行的结果是不一样的，UNION会去除重复数据，UNION ALL不会去除重复数据\n",
	},
	DDLCheckIsExistLimitOffset: {
		Content: "不建议操作\n\n`SELECT column_a FROM table_a LIMIT 5,10`\n\n`SELECT column_a FROM table_a LIMIT 10 OFFSET 5`\n",
	},
	DDLCheckIndexOption: {
		Content: "```\n区分度：可通过字段去重/字段总数统计的方法判断。\n```",
	},
	DDLCheckColumnEnumNotice: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a', \n    column_b ENUM('是','否') DEFAULT NULL COMMENT 'column_b', -- 不建议使用ENUM类型\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckColumnSetNotice: {
		Content: "```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b SET('是','否') DEFAULT NULL COMMENT 'column_b', --不建议使用SET\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckColumnBlobNotice: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a BLOB DEFAULT NULL COMMENT 'column_a',  -- 不建议column使用BLOB类型\n    column_b TEXT DEFAULT NULL COMMENT 'column_b',  -- 不建议column使用TEXT类型\n    PRIMARY KEY (id)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DMLCheckExplainAccessTypeAll: {
		Content: "样例说明：\n```\nEXPLAIN SELECT column_a FROM table_a WHERE ...\n-- 使用EXPLAIN查看扫描行数，是否有索引需要调整\n```\n",
	},
	DMLCheckExplainExtraUsingFilesort: {
		Content: "不建议操作：\n\n```\nSELECT \n  column_a,column_b,column_c\nFROM \n  table_a\nORDER BY\n  column_a,column_b,column_c -- 当MySQL无法使用索引来满足ORDER BY子句的排序要求时，会使用文件排序\n```\n",
	},
	DMLCheckExplainExtraUsingTemporary: {
		Content: "当执行复杂查询或包含临时结果集的查询时，MySQL可能会使用临时表来存储中间结果。这些临时表通常用于以下情况：\n\n```\n1.排序：如果查询包含ORDER BY子句，并且无法使用索引进行排序，MySQL会使用临时表来存储排序的结果。\n2.分组：如果查询包含GROUP BY子句，并且需要计算聚合函数（如SUM、COUNT等），MySQL会使用临时表来存储分组的结果。\n3.连接：如果查询包含多个表的连接操作（如JOIN），MySQL可能会使用临时表来存储连接的中间结果。\n```\n可通过EXPLAIN查看是否使用了临时表",
	},
	DDLCheckCreateView: {
		Content: "样例说明：\n```\nCREATE VIEW   view_example -- 不建议创建视图\nAS \nSELECT column_a, column_b FROM table_a \n```\n",
	},
	DDLCheckCreateTrigger: {
		Content: "样例说明：\n\n```\nCREATE TRIGGER trigger_example\nAFTER INSERT ON table_a\nFOR EACH ROW\nBEGIN -- 触发器内容，违反规则\n    INSERT INTO table_b (column_a) VALUES ('xxx');\nEND\n```\n",
	},
	DDLCheckCreateFunction: {
		Content: "样例说明：\n```\nCREATE FUNCTION custom_function_example()\nRETURNS INT\nBEGIN  --不建议使用自定义函数\n    RETURN 1; \nEND\n```",
	},
	DDLCheckCreateProcedure: {
		Content: "样例说明：\n\n```\nCREATE PROCEDURE procedure_example()\nBEGIN  -- 不建议使用存储过程\n    SELECT column_a FROM table_a;\nEND\n```\n",
	},
	DDLDisableTypeTimestamp: {
		Content: "样例说明\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b TIMESTAMP COMMENT 'column_b',  --禁止使用TIMESTAMP\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DMLCheckAlias: {
		Content: "不建议操作：\n\n```\nSELECT \n  column_a AS column_a  -- 列名a\nFROM \n  table_a AS table_a --表名a\n```\n",
	},
	DDLHintUpdateTableCharsetWillNotUpdateFieldCharset: {
		Content: "不建议修改：\n\n`ALTER TABLE table_a CONVERT TO CHARACTER SET utf8mb4`",
	},
	DDLHintDropColumn: {
		Content: "禁止使用：\n\n```\nALTER TABLE table_a DROP COLUMN column_a\n```\n",
	},
	DDLHintDropPrimaryKey: {
		Content: "禁止使用：\n\n`ALTER TABLE table_a DROP PRIMARY KEY`\n",
	},
	DDLHintDropForeignKey: {
		Content: "禁止使用：\n\n`ALTER TABLE table_a DROP FOREIGN KEY column_f`\n",
	},
	DMLHintInNullOnlyFalse: {
		Content: "样例说明：\n\n```\nSELECT column_a FROM table_a WHERE column_a IN (NULL)  --不建议使用IN (NULL)\n\nSELECT column_a FROM table_a WHERE column_a NOT IN (NULL) --不建议使用NOT IN (NULL)\n```\n",
	},
	DMLCheckSpacesAroundTheString: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a = ' a '  --字符串前后不建议有空格\n```",
	},
	DDLCheckFullWidthQuotationMarks: {
		Content: "样例说明：\n```\nCREATE TABLE “table_a” (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n\n```\nINSERT INTO “table_a”(id) value(1) -- 此时实际创建的表名为“table_a”，而不是table_a\n```",
	},
	DMLNotRecommendOrderByRand: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a \nFROM \n  table_a \nWHERE \n  column_a=0\nORDER BY RAND() --不建议使用 \n```",
	},
	DMLNotRecommendGroupByConstant: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a \nFROM \n  table_a \nWHERE \n  column_a=0 \nGROUP BY 1  --建议使用列名\n```",
	},
	DMLCheckSortDirection: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a=1\nORDER BY column_b,column_c DESC  --ORDER BY中不建议排序\n```",
	},
	DMLHintGroupByRequiresConditions: {
		Content: "样例说明：\n\n```\nSELECT \n  column_b \nFROM \n  table_a \nWHERE \n  column_a=0 \nGROUP BY column_b  \nORDER BY column_b  --GROUP BY语句中建议使用\n```",
	},
	DMLNotRecommendGroupByExpression: {
		Content: "样例说明\n\n```\nSELECT \n  column_a \nFROM \n  table_a \nORDER BY CASE WHEN column_b=3 THEN 1 ELSE column_b END --不建议使用表达式\n```",
	},
	DMLCheckSQLLength: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a,...,column_x  --  注释 ...\nFROM\n  table_a  --  注释 ...\nJOIN  -- 注释 ...\n  ...\nJOIN\n  table_x ON ... AND ...  -- 注释 ...\nWHERE \n  column_a ... AND column_x ...  --  注释 ...\nGROUP BY ...  --  注释 ...\nORDER BY ...  --  注释 ...\n-- 不超过阈值，此值为这条完整的SQL的字符串长度\n```",
	},
	DMLNotRecommendHaving: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a \nFROM \n  table_a \nWHERE \n  column_a=0 \nGROUP BY column_b\nHAVING column_b >1  --不建议条件放HAVING 中\n```",
	},
	DMLHintUseTruncateInsteadOfDelete: {
		Content: "样例说明：\n\n```\nTRUNCATE TABLE table_a  \n```\n```\nDELETE FROM table_a\n```",
	},
	DMLNotRecommendUpdatePK: {
		Content: "样例说明：\n\n```\nUPDATE table_a SET id=2 WHERE id=1  -- 不建议使用，id为主键 \n```",
	},
	DDLCheckColumnQuantity: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_1 INT,...,column_39 INT,  -- column总数不超过阈值，默认值：40\n    PRIMARY KEY (id),\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DDLRecommendTableColumnCharsetSame: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT 'table_a'\n-- 保持字符集一致，列可以不定义字符集，继承使用表的字符集\n```",
	},
	DDLCheckColumnTypeInteger: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT(10) NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b BIGINT(20) NOT NULL DEFAULT 0 COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DDLCheckVarcharSize: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(1024) DEFAULT NULL COMMENT 'column_b',  -- 不超过阈值，默认：1024\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DMLNotRecommendFuncInWhere: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00'    AND column_b+3<>0  --不建议条件中使用函数\n```\n",
	},
	DMLNotRecommendSysdate: {
		Content: "样例说明：\n\n```\nINSERT INTO \n     table_a(column_a) \nVALUES\n    (SYSDATE())   --不建议使用\n```",
	},
	DMLHintSumFuncTips: {
		Content: "样例说明：\n\n```\nSELECT \n SUM(column_b)  --不建议使用SUM(COL)\nFROM \n  table_a \nWHERE \n  column_a=0\n```\n",
	},
	DMLHintCountFuncWithCol: {
		Content: "样例说明：\n\n```\n\nSELECT count(*) FROM table_a WHERE column_a=0  --建议使用 count(*)统计\n\nSELECT count(column_a) FROM table_a WHERE column_a=0  --不建议使用 count(col)\n\n```",
	},
	DDLCheckColumnQuantityInPK: {
		Content: "样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id,column_a)  -- 列数量不超过阈值\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n\n```",
	},
	DMLHintLimitMustBeCombinedWithOrderBy: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a \nFROM \n  table_a \nWHERE \n  column_a=0 \nORDER BY column_a  --建议使用ORDER BY\nLIMIT 100\n```",
	},
	DMLHintTruncateTips: {
		Content: "样例说明：\n\n```\nTRUNCATE TABLE table_a  --不建议使用\n```\n",
	},
	DMLHintDeleteTips: {
		Content: "1.逻辑备份：mysqldump\n\n2.物理备份：XtraBackup\n\n3.增量备份：Binlog",
	},
	DMLCheckSQLInjectionFunc: {
		Content: "样例说明：\n\n```\n禁止使用函数：\n1.sleep\n2.benchmark\n3.get_lock\n4.release_lock\n```",
	},
	DMLCheckNotEqualSymbol: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  column_a<>1 AND column_b !=1  --不建议使用!=，使用<>\n```\n",
	},
	DMLNotRecommendSubquery: {
		Content: "样例说明：\n\n```\nSELECT \n  column_a\nFROM\n  table_a \nWHERE \n  column_a IN (SELECT column_a FROM table_b WHERE column_a>=1)  --不建议使用\n```",
	},
	DDLCheckAutoIncrement: {
		Content: "```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b varchar(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DDLNotAllowRenaming: {
		Content: "禁止使用：\n\n`RENAME TABLE table_a TO table_x`\n\n`ALTER TABLE table_a CHANGE column_a column_x VARCHAR(20)`\n",
	},
	DMLCheckExplainFullIndexScan: {
		Content: "通过EXPLAIN查看是否使用了全索引扫描：\n\n```\ntype字段：该字段表示访问表的方式。如果type的值是ALL，则表示MySQL将执行全表扫描，而不是使用索引进行查询。\nkey字段：该字段显示MySQL选择的索引。如果key的值为NULL，则表示查询将执行全索引扫描。\n```\n\n`如果type是ALL且key是NULL，则很可能发生了全索引扫描。`\n",
	},
	DMLCheckLimitOffsetNum: {
		Content: "样例说明：\n\n`SELECT column_a FROM table_a LIMIT 10 OFFSET 5  -- 不超过阈值`\n",
	},
	DMLCheckUpdateOrDeleteHasWhere: {
		Content: "样例说明：\n\n`UPDATE table_a SET column_a=1 WHERE column_a=0`\n\n`DELETE FROM table_a WHERE column_a=0`",
	},
	DMLCheckSortColumnLength: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    long_text_column VARCHAR(2000) DEFAULT NULL COMMENT 'long_text_column',\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n\nSELECT id FROM table_a ORDER BY long_text_column  -- 不建议使用，对长字段进行排序\n```\n",
	},
	DMLCheckExplainExtraUsingIndexForSkipScan: {
		Content: "当满足以下条件时，MySQL可能会使用索引跳跃扫描：\n\n1.查询条件中涉及到的索引列是联合索引的一部分，但不是第一个列。\n\n2.查询条件中只包含等值匹配，而不是范围查询。\n",
	},
	DMLCheckAffectedRows: {
		Content: "不建议使用\n\n`UPDATE table_a SET column_a=1 WHERE ...  -- 先用SELECT count(*) FROM table_a WHERE ... 查看行数`\n\n`DELETE FROM table_a WHERE ...  -- 先用SELECT count(*) FROM table_a WHERE ... 查看行数`\n\n",
	},
	DMLCheckSameTableJoinedMultipleTimes: {
		Content: "样例说明：\n\n```\nSELECT \n  t_a.column_a\nFROM \n  table_a AS t_a\nJOIN \n  table_b AS t_b1 ON t_a.column_a=t_b1.column_a\nJOIN --不建议单表多次连接\n  table_b AS t_b2 ON t_a.column_a=t_b2.column_b\n```",
	},
	DMLCheckInsertSelect: {
		Content: "INSERT ... SELECT语句可能导致性能问题，特别是在大型数据集上，会导致数据库负载增加和执行时间延长，还可能导致不一致的数据插入从而破坏数据的一致性。\n\n```\n不建议操作：\nINSERT INTO table_a(column_a,column_b,column_c) SELECT column_a,column_b,column_c FROM table_b\n\n```",
	},
	DMLCheckAggregate: {
		Content: "使用聚合函数可能会导致性能问题，特别是在处理大量数据时，会引起不必要的计算开销，影响数据库的查询性能。\n```\n不建议使用：\nSELECT COUNT(*) FROM table_a\nSELECT SUM(column_a) FROM table_a\nSELECT AVG(column_a) FROM table_a\nSELECT MIN(column_a) FROM table_a\nSELECT MAX(column_a) FROM table_a\nSELECT GROUP_CONCAT(column_a) FROM table_a\nSELECT column_a FROM table_a GROUP BY column_a HAVING column_b=0\n```",
	},
	DDLCheckColumnNotNULL: {
		Content: "使用NOT NULL约束可以确保表中的每个记录都包含一个值，有助于维护数据的完整性，使数据的一致性更容易得到保证，数据库优化器也可以更好地执行查询优化，提高查询性能。\n\n样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',  --建议使用NOT NULL约束\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',  --建议使用NOT NULL约束\n    column_b VARCHAR(10) NOT NULL DEFAULT '' COMMENT 'column_b',--建议使用NOT NULL约束\n    PRIMARY KEY (id),  \n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```",
	},
	DMLCheckIndexSelectivity: {
		Content: "如果索引的区分度过小，那么查询优化器可能会选择全表扫描而不是使用索引进行查询，这将导致查询性能下降。因此，为了保证查询效率，建议使用具有较高区分度的索引。区分度：字段去重/字段总数\n\n",
	},
	DDLCheckTableRows: {
		Content: "SELECT count(*) FROM table_a  -- 超过规定值时，根据业务情况，制定清理或归档策略\n",
	},
	DDLCheckCompositeIndexDistinction: {
		Content: "根据最左前缀原则，应该将最常用的查询条件字段放在组合索引的最左侧位置，这样可以最大程度地利用索引的优势，提高查询效率。\n\n样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',  -- 例：选择性第二\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',  -- 例：选择性第一\n    column_c DATE DEFAULT NULL COMMENT 'column_c',  -- 例：选择性第三\n    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',\n    INDEX index_column_b_a_c (column_b, column_a, column_c), \n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n```\n",
	},
	DDLAvoidText: {
		Content: "大字段的存储和查询会占用较多的资源，如果将其与其他字段存放在同一张表中，会导致整张表的性能下降。而将大字段单独存放在一张表中，可以减少对整张表的影响，提高查询效率。\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) NOT NULL DEFAULT '' COMMENT 'column_b',\n    PRIMARY KEY (id),  \n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n\nCREATE TABLE table_b (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    table_a_id INT NOT NULL DEFAULT 0 COMMENT '关联table_a的id字段',  -- 跟主表的主键做关联关系\n    column_t TEXT COMMENT 'column_t',\n    PRIMARY KEY (id),\n    KEY index_a_id (table_a_id)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_b'\n```\n",
	},
	DMLCheckSelectRows: {
		Content: "```\nSELECT count(*) FROM table_a WHERE ...\n-- 使用 count(*)查看数据量，超过10W的，筛选条件必须带上主键或者索引\n\n```",
	},
	DMLCheckScanRows: {
		Content: "```\nEXPLAIN SELECT column_a FROM table_a WHERE ...\n-- 使用EXPLAIN查看扫描行数，超过10W的，筛选条件必须带上主键或者索引\n```\n",
	},
	DMLMustUseLeftMostPrefix: {
		Content: "样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    column_c DATE DEFAULT NULL COMMENT 'column_c',\n    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',\n    INDEX index_column_a_b_c (column_a, column_b, column_c),  --使用联合索引\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n\nSELECT \n  column_a,column_b,column_c,column_d \nFROM\n  table_a \nWHERE \n  column_a=1 AND column_b='a' AND column_c='2020-01-01'; --需遵循最左原则（必须包含column_a条件），否则索引会失效\n```",
	},
	DMLMustMatchLeftMostPrefix: {
		Content: "IN 、OR操作会导致查询无法走全索引，可将SQL拆分为多次等值查询。\n\n样例说明：\n\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    column_c DATE DEFAULT NULL COMMENT 'column_c',\n    column_d DECIMAL(10, 2) DEFAULT 0 COMMENT 'column_d',\n    INDEX index_column_a_b_c (column_a, column_b, column_c),  --建立联合索引\n    PRIMARY KEY (id)\n) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n\nSELECT \n  column_a,column_b,column_c,column_d \nFROM\n  table_a \nWHERE \n  column_a IN (1,2,3) OR column_b='a'  --不建议使用IN 、OR条件进行查询\n```\n",
	},
	DMLCheckJoinFieldUseIndex: {
		Content: "JOIN操作是基于索引的。如果JOIN字段没有索引，那么MySQL需要扫描整个表来找到匹配的行，这会导致查询性能下降。\n\n样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    PRIMARY KEY (id),  \n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_a'\n\nCREATE TABLE table_b (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_b INT NOT NULL DEFAULT 0 COMMENT 'column_b',\n    PRIMARY KEY (id),  \n    KEY index_b (column_b)\n)ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT 'table_b'\n\nSELECT \n  t_a.column_a,t_b.column_b\nFROM \n  table_a AS t_a\nJOIN\n  table_b AS t_b ON t_a.column_a=t_b.column_b  -- JOIN字段需要包含索引\n```\n",
	},
	DMLCheckJoinFieldCharacterSetAndCollation: {
		Content: "索引是按照特定的字符集和排序规则进行存储和排序的，如果关联字段的字符集和排序规则不一致，会导致无法使用索引进行快速查询。\n\n样例说明：\n```\nCREATE TABLE table_a (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_a'  -- 建议定义表级别的字符集和排序规则\n\nCREATE TABLE table_b (\n    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键id',\n    column_a INT NOT NULL DEFAULT 0 COMMENT 'column_a',\n    column_b VARCHAR(10) DEFAULT NULL COMMENT 'column_b',\n    PRIMARY KEY (id),\n    KEY index_a (column_a)\n)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'table_b'  -- 建议定义表级别的字符集和排序规则\n\nSELECT \n  t_a.column_a,t_b.column_b\nFROM \n  table_a AS t_a\nJOIN\n  table_b AS t_b ON t_a.column_a=t_b.column_a \n```\n",
	},
	DMLCheckMathComputationOrFuncOnIndex: {
		Content: "如果对索引列使用了数学运算或函数，会改变其原有的数据结构和排序方式，导致无法使用索引进行快速查询。\n\n样例说明：\n```\nSELECT \n  column_a\nFROM \n  table_a\nWHERE\n  FROM_UNIXTIME(column_a)>'2000-01-01 00:00:00' AND column_b+3>0  -- 例：column_a, column_b为索引列\n```",
	},
}
