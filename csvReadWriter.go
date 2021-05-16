package umscraper

import (
	"encoding/csv"
	"os"
)

// WriteCsv writes a table of [][]string into specified filename.
func WriteCsv(table [][]string, filename string) {
	file, err := os.Create(filename)
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.WriteAll(table)
	checkError("Cannot write file", err)
}

// ReadCsv reads a table of [][]string from specified filename.
func ReadCsv(filename string) [][]string {
	file, err := os.Open(filename)
	checkError("Unable to open file", err)
	defer file.Close()
	table, err := csv.NewReader(file).ReadAll()
	checkError("csv reader Error", err)
	return table
}
