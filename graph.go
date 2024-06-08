package main

import (
	"database/sql"
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"strings"
)

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
		{Key: "penwidth", Value: "1.0"},      // Border width
	}
}

// GenerateERDiagram generates the ER diagram for the specified database
func GenerateERDiagram(username, dbname string) error {
	// Connect to the database
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", username, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Extract tables and columns
	tables, err := extractTables(db)
	if err != nil {
		return err
	}

	// Extract foreign key relationships
	foreignKeys, err := extractForeignKeys(db)
	if err != nil {
		return err
	}

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
		return err
	}

	// Write the DOT data to a file
	dotFileName := "er_diagram.dot"
	err = writeToFile(dotFileName, dotData)
	if err != nil {
		return err
	}

	fmt.Printf("ER diagram saved to %s\n", dotFileName)

	return nil
}
