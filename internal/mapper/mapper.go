package mapper

import (
	"chatgpt-merge/internal/models"
	"log"
	"sort"
	"strconv"
)

func ToSnippets(conversations []models.Conversation) []models.Snippet {
	var snippets []models.Snippet

	for _, conversation := range conversations {
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

func ToCSVRow(snippet models.Snippet) ([]string, error) {
	createTimeStr := strconv.FormatFloat(snippet.CreateTime, 'f', -1, 64)

	return []string{createTimeStr, snippet.Role, snippet.Content}, nil
}
