package mapper

import (
	"chatgpt-merge/internal/models"
	"sort"
	"strings"
	"time"
)

func MapToSnippets(conversations []models.Conversation, includeTitles []string) []models.Snippet {
	snippets := []models.Snippet{}

	// Build a set of included titles after trimming; skip empty lines
	titleSet := make(map[string]struct{}, len(includeTitles))
	for _, t := range includeTitles {
		st := strings.TrimSpace(t)
		if st == "" {
			continue
		}
		titleSet[st] = struct{}{}
	}

	for _, conversation := range conversations {
		// ignore conversations that are not selected
		if _, ok := titleSet[conversation.Title]; !ok {
			continue
		}

		for _, mapping := range conversation.Mapping {
			msg := mapping.Message
			if msg.CreateTime == 0 {
				continue
			}
			if msg.Author.Role != "assistant" && msg.Author.Role != "user" {
				continue
			}
			parts := msg.Content.Parts
			if len(parts) == 0 {
				continue
			}
			// Collect only string parts; skip non-strings
			var sb []string
			for _, p := range parts {
				if s, ok := p.(string); ok && s != "" {
					sb = append(sb, s)
				}
			}
			if len(sb) == 0 {
				continue
			}

			snippets = append(snippets, models.Snippet{
				CreateTime: time.Unix(int64(msg.CreateTime), 0).UTC().Format(time.RFC3339),
				Content:    strings.Join(sb, "\n"),
				Role:       msg.Author.Role,
			})
		}
	}

	// Sort by timestamp (RFC3339 strings sort lexicographically by time, but be explicit by parsing)
	sort.Slice(snippets, func(i, j int) bool {
		// RFC3339 lexical order equals chronological order, so compare strings
		return snippets[i].CreateTime < snippets[j].CreateTime
	})

	return snippets
}

func MapToCSVRow(snippet models.Snippet) []string {
	return []string{snippet.CreateTime, snippet.Role, snippet.Content}
}
