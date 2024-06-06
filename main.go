package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
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

// LabeledNode is a custom node type that includes a label.
type LabeledNode struct {
	graph.Node
	Label string
}

// Attributes returns the attributes for the node in DOT format.
func (n LabeledNode) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "label", Value: n.Label},
		{Key: "shape", Value: "box"},         // Custom shape
		{Key: "style", Value: "filled"},      // Filled style
		{Key: "color", Value: "black"},       // Border color
		{Key: "fillcolor", Value: "#D3D3D3"}, // Fill color
		{Key: "fontname", Value: "Arial"},    // Font name
		{Key: "fontsize", Value: "12"},       // Font size
	}
}

func main() {
	// Get database connection parameters from user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter database username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter database name: ")
	dbname, _ := reader.ReadString('\n')
	dbname = strings.TrimSpace(dbname)

	// Connect to the database
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", username, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Extract tables and columns
	tables := extractTables(db)

	// Extract foreign key relationships
	foreignKeys := extractForeignKeys(db)

	// Create the graph
	g := simple.NewDirectedGraph()
	nodes := make(map[string]graph.Node)

	// Add nodes to the graph
	for _, table := range tables {
		label := fmt.Sprintf("%s.%s\n%s", table.Schema, table.Name, strings.Join(table.Columns, ", "))
		node := LabeledNode{
			Node:  g.NewNode(),
			Label: label,
		}
		g.AddNode(node)
		nodes[fmt.Sprintf("%s.%s", table.Schema, table.Name)] = node
	}

	// Add edges (foreign keys) to the graph
	for _, fk := range foreignKeys {
		from := nodes[fmt.Sprintf("%s.%s", fk.Schema, fk.TableName)]
		to := nodes[fmt.Sprintf("%s.%s", fk.ForeignSchema, fk.ForeignTableName)]
		g.SetEdge(g.NewEdge(from, to))
	}

	// Export the graph to DOT format
	dotData, err := dot.Marshal(g, "ER Diagram", "", "")
	if err != nil {
		log.Fatal(err)
	}

	// Write the DOT data to a file
	dotFileName := "er_diagram.dot"
	f, err := os.Create(dotFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(dotData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ER diagram saved to %s\n", dotFileName)

	// Run Graphviz `dot` command to generate PNG
	pngFileName := "er_diagram.png"
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", pngFileName)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ER diagram image saved to %s\n", pngFileName)
}

func extractTables(db *sql.DB) []Table {
	rows, err := db.Query(`
        SELECT table_schema, table_name, column_name
        FROM information_schema.columns
        WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tableMap := make(map[string]*Table)
	for rows.Next() {
		var schema, tableName, columnName string
		if err := rows.Scan(&schema, &tableName, &columnName); err != nil {
			log.Fatal(err)
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
	return tables
}

func extractForeignKeys(db *sql.DB) []ForeignKey {
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
		log.Fatal(err)
	}
	defer rows.Close()

	var foreignKeys []ForeignKey
	for rows.Next() {
		var fk ForeignKey
		if err := rows.Scan(&fk.Schema, &fk.TableName, &fk.ColumnName,
			&fk.ForeignSchema, &fk.ForeignTableName, &fk.ForeignColumnName); err != nil {
			log.Fatal(err)
		}
		foreignKeys = append(foreignKeys, fk)
	}
	return foreignKeys
}
