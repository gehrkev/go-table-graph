package main

import (
	"database/sql"
	"fmt"
)

// Table represents a database table with its columns.
type Table struct {
	Schema  string
	Name    string
	Columns []string
}

// ForeignKey represents a foreign key relationship between tables.
type ForeignKey struct {
	Schema            string
	TableName         string
	ColumnName        string
	ForeignSchema     string
	ForeignTableName  string
	ForeignColumnName string
}

// extractTables extracts tables and their columns from the database
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

// extractForeignKeys extracts foreign key relationships from the database
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
