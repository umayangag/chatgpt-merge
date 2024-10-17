package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	maxFileSize = 16 * 1024 * 1024 // 16 MB
)

// Structure for your JSON object (modify it according to your actual structure)
type Conversation struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	User    string `json:"user"`
}

func splitJSONFile(filename string) error {
	// Open the original file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the JSON array from the file
	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Initialize variables for splitting
	var part int
	var currentSize int64
	var splitConversations []Conversation

	// Iterate through the conversations to split into multiple files
	for _, conversation := range conversations {
		// Estimate the size of the current conversation (as JSON) to check the file size limit
		conversationBytes, err := json.Marshal(conversation)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		// If adding this object exceeds the max file size, write the current file
		if currentSize+int64(len(conversationBytes)) > maxFileSize {
			err = writePartFile(filename, part, splitConversations)
			if err != nil {
				return err
			}

			// Start a new part file
			part++
			splitConversations = []Conversation{} // Reset the array for the new part
			currentSize = 0
		}

		// Add the conversation to the current part
		splitConversations = append(splitConversations, conversation)
		currentSize += int64(len(conversationBytes))
	}

	// Write the remaining conversations if any
	if len(splitConversations) > 0 {
		err = writePartFile(filename, part, splitConversations)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to write part files
func writePartFile(filename string, part int, conversations []Conversation) error {
	// Create the part file
	partFilename := fmt.Sprintf("%s.part%d.json", filename, part)
	partFile, err := os.Create(partFilename)
	if err != nil {
		return fmt.Errorf("failed to create part file: %w", err)
	}
	defer partFile.Close()

	// Write the JSON array to the file
	encoder := json.NewEncoder(partFile)
	encoder.SetIndent("", "  ") // Optional: for pretty-printing
	err = encoder.Encode(conversations)
	if err != nil {
		return fmt.Errorf("failed to encode JSON to part file: %w", err)
	}

	fmt.Printf("Created part file: %s\n", partFilename)
	return nil
}

func main() {
	// Example JSON file path
	filename := "conversations.json"

	err := splitJSONFile(filename)
	if err != nil {
		fmt.Printf("Error splitting file: %v\n", err)
	} else {
		fmt.Println("File split successfully!")
	}
}
