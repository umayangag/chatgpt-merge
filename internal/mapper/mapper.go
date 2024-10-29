package mapper

import (
	"chatgpt-merge/internal/models"
	"log"
	"sort"
	"strconv"

	"golang.org/x/exp/slices"
)

var ignoreTitles = []string{}

func MapToSnippets(conversations []models.Conversation) []models.Snippet {
	snippets := []models.Snippet{}

	for _, conversation := range conversations {
		if slices.Contains(ignoreTitles, conversation.Title) {
			log.Println("Ignoring Conversation:", conversation.Title)
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

	sort.Slice(snippets, func(i, j int) bool {
		return snippets[i].CreateTime < snippets[j].CreateTime
	})

	return snippets
}

func MapToCSVRow(snippet models.Snippet) []string {
	createTimeStr := strconv.FormatFloat(snippet.CreateTime, 'f', -1, 64)

	return []string{createTimeStr, snippet.Role, snippet.Content}
}
