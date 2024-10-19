package main

import (
	"chatgpt-merge/internal/models"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Read the JSON file
	data, err := os.ReadFile("input/conversations.json")
	if err != nil {
		log.Println("Error reading the file:", err)
		return
	}

	// Unmarshal the JSON data into an array of items
	var conversations []models.Conversation
	err = json.Unmarshal(data, &conversations)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	snippets := generateSnippets(conversations)

	log.Println("No. of snippets extracted:", len(snippets))

	err = writeSnippetsToCSV("output/output.csv", snippets)
	if err != nil {
		log.Println("Error writing csv file:", err)
		return
	}

	log.Println("Data successfully extracted")
}

func generateSnippets(conversations []models.Conversation) []models.Snippet {
	// Prepare a slice to hold the output data
	var snippets []models.Snippet

	// Extract create_time, content and role from each message
	for _, conversation := range conversations {
		if conversation.Title == "ChatGPT Conversation Merging Script" {
			continue
		}
		log.Println("Scanning Conversation:", conversation.Title)

		for _, mapping := range conversation.Mapping {
			msg := mapping.Message
			content := msg.Content.Parts
			if msg.CreateTime == 0 {
				continue
			}
			if msg.Author.Role != "assistant" && msg.Author.Role != "user" {
				continue
			}
			if len(content) == 0 {
				continue
			}
			context, ok := content[0].(string)
			if !ok {
				continue
			}

			snippets = append(snippets, models.Snippet{
				CreateTime: msg.CreateTime,
				Content:    context,
				Role:       msg.Author.Role,
			})
		}
	}

	// Sort the output by create_time in ascending order
	sort.Slice(snippets, func(i, j int) bool {
		return snippets[i].CreateTime < snippets[j].CreateTime
	})

	return snippets
}

// Write an array of Snippet structs to a CSV file
func writeSnippetsToCSV(filename string, snippets []models.Snippet) error {
	// Create a CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Timestamp", "Role", "Content"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write each Snippet as a row
	for _, snippet := range snippets {
		row, err := snippetToCSVRow(snippet)
		if err != nil {
			return err
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Convert the Snippet struct to a CSV row
func snippetToCSVRow(snippet models.Snippet) ([]string, error) {
	// Convert CreateTime (float64) to string
	createTimeStr := strconv.FormatFloat(snippet.CreateTime, 'f', -1, 64)

	// Return a CSV row
	return []string{createTimeStr, snippet.Role, snippet.Content}, nil
}
