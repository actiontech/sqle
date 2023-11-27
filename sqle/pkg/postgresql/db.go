package postgresql

import "C"
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type DSN struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (d *DSN) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Database)
}

type DB struct {
	Db              *sql.DB
	IsCaseSensitive bool
}

func NewDB(dsn *DSN) (*DB, error) {
	// 创建一个数据库连接池
	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return nil, err
	}

	// 设置连接池的最大连接数和空闲连接数
	db.SetMaxOpenConns(100) // 设置最大连接数
	db.SetMaxIdleConns(10)  // 设置空闲连接数

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{Db: db}, nil
}

func (o *DB) Close() error {
	return o.Db.Close()
}

func (o *DB) GetCaseSensitive() bool {
	var isCaseSensitive bool
	query := "SELECT setting FROM pg_settings WHERE name = 'quote_all_identifiers'"

	sqls, err := getResultSqls(o.Db, query)
	if err != nil {
		return false
	}
	if len(sqls) == 0 {
		return false
	}
	for _, sqlContent := range sqls {
		if strings.ToLower(sqlContent) == "on" {
			return true
		}
	}
	return isCaseSensitive
}

func (o *DB) GetAllUserSchemas() ([]string, error) {
	query := "SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' AND nspname != 'information_schema'"

	sqls, err := getResultSqls(o.Db, query)
	if err != nil {
		return nil, err
	}
	return sqls, nil
}

func (o *DB) ShowSchemaTables(schema string) ([]string, error) {
	query := fmt.Sprintf("select TABLE_NAME from information_schema.tables "+
		" where table_schema='%s' and TABLE_TYPE in ('BASE TABLE','SYSTEM VIEW')", schema)

	if o.IsCaseSensitive {
		schema = strings.ToLower(schema)
		query = fmt.Sprintf("select TABLE_NAME from information_schema.tables "+
			" where lower(table_schema)='%s' and TABLE_TYPE in ('BASE TABLE','SYSTEM VIEW')", schema)
	}
	return getResultSqls(o.Db, query)
}

func (o *DB) ShowSchemaViews(schema string) ([]string, error) {
	query := fmt.Sprintf("select TABLE_NAME from information_schema.tables "+
		" where table_schema='%s' and TABLE_TYPE='VIEW'", schema)

	if o.IsCaseSensitive {
		schema = strings.ToLower(schema)
		query = fmt.Sprintf("select TABLE_NAME from information_schema.tables "+
			"where lower(table_schema)='%s' and TABLE_TYPE='VIEW'", schema)
	}
	return getResultSqls(o.Db, query)
}

func (o *DB) ShowCreateTables(database, schema, tableName string) ([]string, error) {
	tables := make([]string, 0)
	tableDDl := fmt.Sprintf("CREATE TABLE %s.%s(", schema, tableName)
	if o.IsCaseSensitive {
		database = strings.ToLower(database)
		schema = strings.ToLower(schema)
		tableName = strings.ToLower(tableName)
	}
	columnsCondition := fmt.Sprintf("table_catalog = '%s' AND table_schema = '%s' AND table_name = '%s'",
		database, schema, tableName)
	if o.IsCaseSensitive {
		columnsCondition = fmt.Sprintf("lower(table_catalog) = '%s' AND lower(table_schema) = '%s' "+
			"AND lower(table_name) = '%s'", database, schema, tableName)
	}
	// 获取列定义，多个英文逗号分割
	columns := fmt.Sprintf("SELECT string_agg(column_name || ' ' || "+
		"CASE "+
		" WHEN data_type IN ('char', 'varchar', 'character', 'character varying', 'text') "+
		" THEN data_type || '(' || COALESCE(character_maximum_length, 0) || ')' "+
		" WHEN data_type IN ('numeric', 'decimal') "+
		" THEN data_type || '(' || COALESCE(numeric_precision, 0) || ',' || COALESCE(numeric_scale, 0) || ')' "+
		" WHEN data_type IN ('integer', 'smallint', 'bigint') THEN data_type "+
		" ELSE data_type "+
		" END "+
		" || "+
		" CASE "+
		" WHEN column_default != '' THEN ' DEFAULT ' || column_default ELSE '' END "+
		" || "+
		" CASE "+
		" WHEN is_nullable = 'NO' THEN ' NOT NULL' ELSE '' END, ',\n ' ORDER BY ordinal_position) AS columns_sql"+
		" FROM information_schema.columns "+
		" WHERE %s GROUP BY table_name", columnsCondition)
	sqls, err := getResultSqls(o.Db, columns)
	if err != nil {
		log.Printf("search column definition error:%s\n", err)
		return nil, err
	}
	if len(sqls) == 0 {
		return tables, nil
	}
	tableDDl += strings.Join(sqls, "")
	constraintsCondition := fmt.Sprintf("n.nspname = '%s' AND C.relname = '%s'", schema, tableName)
	if o.IsCaseSensitive {
		constraintsCondition = fmt.Sprintf("lower(n.nspname) = '%s' "+
			"AND lower(C.relname) = '%s'", schema, tableName)
	}
	// 获取所有约束
	constraints := fmt.Sprintf("SELECT 'CONSTRAINT ' || r.conname || ' ' || "+
		" pg_catalog.pg_get_constraintdef ( r.OID, TRUE ) AS constraint_definition "+
		" FROM pg_catalog.pg_constraint r "+
		" JOIN pg_catalog.pg_class C ON C.OID = r.conrelid "+
		" JOIN pg_catalog.pg_namespace n ON n.OID = C.relnamespace "+
		" WHERE %s", constraintsCondition)
	sqls, err = getResultSqls(o.Db, constraints)
	if err != nil {
		log.Printf("search constraint definition error:%s\n", err)
		return nil, err
	}
	for _, sqlContext := range sqls {
		tableDDl += ",\n" + sqlContext
	}
	tableDDl += ")"
	indexesCondition := fmt.Sprintf("schemaname = '%s' and tablename = '%s' ", schema, tableName)
	if o.IsCaseSensitive {
		indexesCondition = fmt.Sprintf("lower(schemaname) = '%s' and lower(tablename) = '%s'",
			schema, tableName)
	}
	// 获取索引
	indexes := fmt.Sprintf("SELECT indexdef AS index_definition FROM pg_indexes "+
		" WHERE %s", indexesCondition)
	sqls, err = getResultSqls(o.Db, indexes)
	if err != nil {
		log.Printf("search index definition error:%s\n", err)
		return nil, err
	}
	for _, sqlContent := range sqls {
		if strings.Contains(sqlContent, "CREATE UNIQUE INDEX") {
			continue
		}
		tableDDl += ";\n" + sqlContent
	}
	tables = append(tables, tableDDl)
	return tables, nil
}

func (o *DB) ShowCreateViews(database, schema, tableName string) ([]string, error) {
	query := fmt.Sprintf(
		"SELECT 'CREATE OR REPLACE VIEW ' || table_schema || '.' || table_name || ' AS ' || view_definition"+
			" AS create_view_statement "+
			" FROM information_schema.views "+
			" WHERE table_catalog = '%s' AND table_schema = '%s' AND table_name = '%s'",
		database, schema, tableName)

	if o.IsCaseSensitive {
		database = strings.ToLower(database)
		tableName = strings.ToLower(tableName)
		query = fmt.Sprintf(
			"SELECT 'CREATE OR REPLACE VIEW ' || table_schema || '.' || table_name || ' AS ' || view_definition"+
				" AS create_view_statement "+
				" FROM information_schema.views "+
				" WHERE lower(table_catalog) = '%s' AND lower(table_schema) = '%s' AND lower(table_name) = '%s'",
			database, schema, tableName)
	}
	return getResultSqls(o.Db, query)
}

func getResultSqls(db *sql.DB, query string) ([]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		innerErr := rows.Close()
		if innerErr != nil {
			log.Printf("Close rows error:%s\n", innerErr)
		}
	}(rows)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	sqls := make([]string, 0)
	for rows.Next() {
		var sqlContent string
		err = rows.Scan(&sqlContent)
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sqlContent)
	}
	return sqls, nil
}
