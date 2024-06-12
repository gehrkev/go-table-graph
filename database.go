package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Table represents a database table with its columns.
type Table struct {
	Schema  string   `json:"schema"`
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

// ForeignKey represents a foreign key relationship between tables.
type ForeignKey struct {
	Schema            string `json:"schema"`
	TableName         string `json:"tableName"`
	ColumnName        string `json:"columnName"`
	ForeignSchema     string `json:"foreignSchema"`
	ForeignTableName  string `json:"foreignTableName"`
	ForeignColumnName string `json:"foreignColumnName"`
}

// extractTables extracts tables and their columns from the database.
// It queries the information_schema.columns table to retrieve table names,
// column names, and their respective schemas, excluding system schemas.
func extractTables(db *sql.DB) ([]Table, error) {
	rows, err := db.Query(`
		SELECT table_schema, table_name, column_name
		FROM information_schema.columns
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tableMap := make(map[string]*Table)
	for rows.Next() {
		var schema, tableName, columnName string
		if err := rows.Scan(&schema, &tableName, &columnName); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%s.%s", schema, tableName)
		if tableMap[key] == nil {
			tableMap[key] = &Table{
				Schema: schema,
				Name:   tableName,
			}
		}
		tableMap[key].Columns = append(tableMap[key].Columns, columnName)
	}

	var tables []Table
	for _, table := range tableMap {
		tables = append(tables, *table)
	}
	return tables, nil
}

// extractForeignKeys extracts foreign key relationships from the database.
// It queries the information_schema tables to retrieve foreign key constraints
// and their related tables and columns.
func extractForeignKeys(db *sql.DB) ([]ForeignKey, error) {
	rows, err := db.Query(`
		SELECT
			tc.table_schema,
			tc.table_name,
			kcu.column_name,
			ccu.table_schema AS foreign_table_schema,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name
		FROM
			information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
			JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []ForeignKey
	for rows.Next() {
		var fk ForeignKey
		if err := rows.Scan(&fk.Schema, &fk.TableName, &fk.ColumnName,
			&fk.ForeignSchema, &fk.ForeignTableName, &fk.ForeignColumnName); err != nil {
			return nil, err
		}
		foreignKeys = append(foreignKeys, fk)
	}
	return foreignKeys, nil
}

// ExtractERDiagram connects to the database using provided username and dbname,
// extracts tables and foreign key relationships, and returns them.
func ExtractERDiagram(username, password, dbname string) ([]Table, []ForeignKey, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)
	fmt.Println("Connecting to database with:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	tables, err := extractTables(db)
	if err != nil {
		return nil, nil, fmt.Errorf("error extracting tables: %v", err)
	}

	foreignKeys, err := extractForeignKeys(db)
	if err != nil {
		return nil, nil, fmt.Errorf("error extracting foreign keys: %v", err)
	}

	return tables, foreignKeys, nil
}
