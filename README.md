
# go-table-graph

## Overview

````go-table-graph```` is a web application written in Go that allows you to visualize database table relationships as an Entity-Relationship (ER) diagram. It connects to local PostgreSQL databases, extracts table structures, and renders the ER diagram using the vis-network JavaScript library for interactive visualization.

## Requirements

- **Go**: Ensure that you have Go installed on your system. You can download it from [go.dev](https://go.dev/dl/).

- **PostgreSQL Database**: You need access to a local PostgreSQL database to connect to and extract table relationships.

## Installation and Setup

### Clone the Repository

```bash
git clone https://github.com/gehrkev/go-table-graph.git
cd go-table-graph
```

### Build and Run

1. Build the executable:

   ```bash
   go build -o go-table-graph
   ```

2. Run the application:

   ```bash
   ./go-table-graph
   ```

3. Open your web browser and go to `http://localhost:8080` to access the application.

## Usage

1. Enter your PostgreSQL database username, password and database name in the input fields.
2. Click on **Generate ER Diagram** to fetch and visualize the database schema.


## External Libraries Used

- **[gin](https://github.com/gin-gonic/gin)**: Gin is a HTTP web framework written in Go (Golang).
- **[pq](https://github.com/lib/pq)**: Pure Go Postgres driver for database/sql.
- **[vis-network](https://github.com/visjs/vis-network)**: 💫 Display dynamic, automatically organised, customizable network views.

