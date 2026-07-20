package main

import (
	"encoding/json"
	"fmt"

	"github.com/peterjohnbishop/friendly-octo-enigma/models"
	"github.com/peterjohnbishop/friendly-octo-enigma/processors"
	"github.com/peterjohnbishop/friendly-octo-enigma/server"
	"github.com/peterjohnbishop/friendly-octo-enigma/tui"
)

func main() {
	rawHeadersChan := make(chan map[string][]string, 100)
	rawBodyChan := make(chan []byte, 100)

	outputHeadersChan := make(chan map[string][]string, 100)
	outputBodyChan := make(chan map[string]models.PropertyDetail, 100)

	// the processors startup and wait for raw input, then output the result to the channels to the TUI
	go processors.MapAndMergeHeaders(rawHeadersChan, outputHeadersChan)
	go processors.MapBody(rawBodyChan, outputBodyChan)

	// the server hosts the webhook endpoint where raw input is sent, which is output to the raw channels
	go server.ServeGin(rawHeadersChan, rawBodyChan)

	// Start TUI on the main thread
	if err := tui.Start(outputHeadersChan, outputBodyChan); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
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
