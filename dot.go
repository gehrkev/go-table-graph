package main

import (
	"fmt"
	"os"
	"os/exec"
)

// writeToFile writes data to a file
func writeToFile(fileName string, data []byte) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", fileName, err)
	}
	return nil
}

// ConvertDotToPNG converts a DOT file to PNG using the dot tool
func ConvertDotToPNG(dotFileName, pngFileName string) error {
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", pngFileName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute dot command: %v", err)
	}
	fmt.Printf("Converted %s to %s successfully\n", dotFileName, pngFileName)
	return nil
}
