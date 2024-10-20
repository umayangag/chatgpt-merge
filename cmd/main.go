package main

import (
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"chatgpt-merge/internal/writer"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Invalid arguments. Source path and output path are required")
		os.Exit(1)
	}

	source := args[0]
	output := args[1]

	data, err := os.ReadFile(source)
	if err != nil {
		log.Fatal("Error reading the file:", err)
		return
	}

	var conversations []models.Conversation
	err = json.Unmarshal(data, &conversations)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
		return
	}

	snippets := mapper.MapToSnippets(conversations)

	log.Println("No. of snippets extracted:", len(snippets))

	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Error creating output file:", err)
		return
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	err = writer.WriteToCSV(csvWriter, mapper.MapToCSVRow, snippets)
	if err != nil {
		log.Fatal("Error writing csv file:", err)
		return
	}

	log.Println("Data successfully extracted")
}
