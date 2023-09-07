package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ffelipelimao/compiler/compiler"
)

type LogEntry struct {
	Time       string            `json:"time"`
	Event      string            `json:"event"`
	Attributes map[string]string `json:"attributes"`
}

func main() {

	c := compiler.New("qgames.log")

	if err := c.LoadRows(); err != nil {
		log.Panic("error to read rows", err)
	}

	event := c.Process()

	rankingsJson, _ := json.Marshal(event)
	os.WriteFile("output.json", rankingsJson, 0644)

}
