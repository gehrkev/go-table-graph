package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.StaticFile("/", "./index.html")

	r.GET("/api/erdiagram", handleERDiagramRequest)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Unable to start server:", err)
	}
}

// handleERDiagramRequest handles GET requests to "/api/erdiagram".
// It retrieves database username and dbname from query parameters,
// extracts ER diagram data using ExtractERDiagram function, and returns it as JSON.
func handleERDiagramRequest(c *gin.Context) {
	username := c.Query("username")
	dbname := c.Query("dbname")
	if username == "" || dbname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and dbname are required"})
		return
	}

	tables, foreignKeys, err := ExtractERDiagram(username, dbname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tables":      tables,
		"foreignKeys": foreignKeys,
	})
}
