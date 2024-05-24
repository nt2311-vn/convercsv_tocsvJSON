package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir := "./data"

	csvFile, err := findCSVFile(dir)
	if err != nil {
		log.Fatalf("Error finding CSV file: %v", err)
	}

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error on read all records: %s", err)
	}

	header := records[0]
	data := records[1:]

	groupedData := make(map[string]*GroupRecords)

	for _, row := range data {
		// print out row
		for i, cell := range row {
			fmt.Printf("%s: %s\n", header[i], cell)
		}
	}
}

func findCSVFile(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no CSV file found in %s", dir)
}
