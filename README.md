# go-table-graph

## Overview

`go-table-graph` is a Go CLI tool to generate and visualize database table relationships as an ER diagram using Graphviz's DOT tool. It connects to PostgreSQL databases, extracts table structures and foreign key relationships, and outputs an ER diagram in PNG format.

## Requirements

- **Go Programming Language**: Make sure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

- **Graphviz**: You need Graphviz installed on your system to convert the DOT file to PNG format. Graphviz can be downloaded from [graphviz.org](https://graphviz.org/download/).

## Usage

### Running with `go run`

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go-table-graph.git
   cd go-table-graph
   ```

2. Execute the program using `go run`:

   ```bash
   go run main.go database.go graph.go dot.go
   ```

   Follow the prompts to enter your database username and database name.

### Compiling with `go build`

1. Clone the repository (if not already cloned):

   ```bash
   git clone https://github.com/yourusername/go-table-graph.git
   cd go-table-graph
   ```

2. Compile the program:

   ```bash
   go build -o go-table-graph
   ```

   This will create an executable named `go-table-graph` in the current directory.

4. Run the compiled executable:

   ```bash
   ./go-table-graph
   ```

   Follow the prompts to enter your database username and database name.

### Example

After running the program, it will generate an ER diagram in PNG format (`er_diagram.png`) based on your PostgreSQL database schema.
