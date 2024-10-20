package main

import (
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"chatgpt-merge/internal/writer"
	"encoding/json"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Println("invalid arguments. source path and output path are required")
		os.Exit(1)
	}

	source := args[0]
	output := args[1]

	data, err := os.ReadFile(source)
	if err != nil {
		log.Println("Error reading the file:", err)
		return
	}

	var conversations []models.Conversation
	err = json.Unmarshal(data, &conversations)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	snippets := mapper.ToSnippets(conversations)

	log.Println("No. of snippets extracted:", len(snippets))

	err = writer.ToCSV(output, snippets)
	if err != nil {
		log.Println("Error writing csv file:", err)
		return
	}

	log.Println("Data successfully extracted")
}
