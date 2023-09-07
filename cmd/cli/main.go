package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ffelipelimao/compiler/internal/processor"
	"github.com/ffelipelimao/compiler/internal/reader"
)

func main() {
	reader := reader.New("./qgames.log")

	rows, err := reader.LoadRows()
	if err != nil {
		log.Panic("error to load .log file", err)
	}

	if rows == nil {
		log.Panic("error to read the rows")
	}

	processor := processor.New()
	output := processor.CreateOutput(rows)

	outputJSON, err := json.Marshal(output)
	if err != nil {
		log.Panic("error to marshall output", err)
	}

	os.WriteFile("output.json", outputJSON, 0644)
}
