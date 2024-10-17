package main

import (
	"chatgpt-merge/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func main() {
	// Read the JSON file
	data, err := os.ReadFile("conversations.json")
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	// Unmarshal the JSON data into an array of items
	var conversations []models.Conversation
	err = json.Unmarshal(data, &conversations)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Prepare a slice to hold the output data
	var snippets []models.Snippet

	// Extract create_time, content and role from each message
	for _, conversation := range conversations {
		for _, mapping := range conversation.Mapping {
			msg := mapping.Message
			content := msg.Content.Parts
			if msg.CreateTime != 0 && content != nil {
				snippets = append(snippets, models.Snippet{
					CreateTime: msg.CreateTime,
					Content:    content,
					Role:       msg.Author.Role,
				})
			}
		}
	}

	// Sort the output by create_time in ascending order
	sort.Slice(snippets, func(i, j int) bool {
		return snippets[i].CreateTime < snippets[j].CreateTime
	})

	// Marshal the output data to JSON
	jsonExport, err := json.Marshal(snippets)
	if err != nil {
		fmt.Println("Error marshalling output to JSON:", err)
		return
	}

	// Write the output to a new JSON file
	err = os.WriteFile("output.json", jsonExport, 0644)
	if err != nil {
		fmt.Println("Error writing the output file:", err)
		return
	}

	fmt.Println("Data successfully extracted to output.json")
}
