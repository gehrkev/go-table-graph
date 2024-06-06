package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter database username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter database name: ")
	dbname, _ := reader.ReadString('\n')
	dbname = strings.TrimSpace(dbname)

	err := GenerateERDiagram(username, dbname)
	if err != nil {
		log.Fatal(err)
	}

	dotFileName := "er_diagram.dot"
	pngFileName := "er_diagram.png"
	err = ConvertDotToPNG(dotFileName, pngFileName)
	if err != nil {
		log.Fatal(err)
	}
}
