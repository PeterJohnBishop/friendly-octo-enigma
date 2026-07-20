package main

import (
	"encoding/json"
	"fmt"

	"github.com/peterjohnbishop/friendly-octo-enigma/models"
	"github.com/peterjohnbishop/friendly-octo-enigma/processors"
	"github.com/peterjohnbishop/friendly-octo-enigma/server"
)

func main() {
	rawHeaders := make(chan map[string][]string, 100)
	rawBody := make(chan []byte, 100)

	outputHeadersChan := make(chan map[string][]string, 100)
	outputBodyChan := make(chan map[string]models.PropertyDetail, 100)

	go processors.MapAndMergeHeaders(rawHeaders, outputHeadersChan)
	go processors.MapBody(rawBody, outputBodyChan)
	go server.ServeGin(rawHeaders, rawBody)

	for {
		select {
		case header := <-outputHeadersChan:
			prettyPrint("--- HEADER ---", header)

		case body := <-outputBodyChan:
			prettyPrint("--- BODY ---", body)
		}
	}
}

func prettyPrint(title string, data any) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("%s\nError formatting JSON: %v\n\n", title, err)
		return
	}

	fmt.Printf("%s\n%s\n\n", title, string(prettyJSON))
}
