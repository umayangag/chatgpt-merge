package mapper

import (
	"chatgpt-merge/internal/models"
	"log"
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

func MapToSnippets(conversations []models.Conversation, includeTitles []string) []models.Snippet {
	snippets := []models.Snippet{}

	for _, conversation := range conversations {
		// ignore conversations that are not selected
		if !slices.Contains(includeTitles, conversation.Title) {
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
				CreateTime: time.Unix(int64(msg.CreateTime), 0).UTC().String(),
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

	return []string{snippet.CreateTime, snippet.Role, snippet.Content}
}
