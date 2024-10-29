package main

import (
	"bufio"
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"chatgpt-merge/internal/writer"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Parse Arguments
	includeFrom := flag.String("include", "", "filepath to list of conversation files to be selected for merging")
	isDryRun := flag.Bool("dry", false, "output the list conversations without merging")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Invalid arguments. Source path and output path are required")
		os.Exit(1)
	}

	source := args[0]
	output := args[1]

	// Read Source File
	data, err := os.ReadFile(source)
	if err != nil {
		log.Fatal("Error reading the source file:", err)
		os.Exit(1)
	}

	var conversations []models.Conversation
	err = json.Unmarshal(data, &conversations)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	if *isDryRun {
		err = listConversations(conversations, output)
		if err != nil {
			log.Fatal("Error writing output file:", err)
			os.Exit(1)
		}
		log.Println("Conversation List successfully extracted to:", output)
		os.Exit(0)
	}

	if includeFrom == nil || *includeFrom == "" {
		log.Fatal("path to list of conversations is required. use flag -include or use -dry to output the list of conversations")
		os.Exit(1)
	}

	titlesData, err := os.ReadFile(*includeFrom)
	if err != nil {
		log.Fatal("Error reading the file for selected conversations:", err)
		os.Exit(1)
	}

	includeTitles := strings.Split(string(titlesData), "\n")
	snippets := mapper.MapToSnippets(conversations, includeTitles)

	log.Println("No. of snippets extracted:", len(snippets))

	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Error creating output file:", err)
		os.Exit(1)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	err = writer.WriteToCSV(csvWriter, mapper.MapToCSVRow, snippets)
	if err != nil {
		log.Fatal("Error writing csv file:", err)
		os.Exit(1)
	}

	log.Println("Data successfully extracted to:", output)
}

func listConversations(conversations []models.Conversation, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(file)

	log.Println("Following Conversations were detected:")
	for _, conversation := range conversations {
		_, err := writer.WriteString(conversation.Title + "\n")
		if err != nil {
			return err
		}
		fmt.Println(conversation.Title)
	}
	return writer.Flush()
}
