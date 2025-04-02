package mapper_test

import (
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapToSnippets(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		input    []models.Conversation
		titles   []string
		expected []models.Snippet
	}{
		{
			name: "valid conversation",
			input: []models.Conversation{
				{
					Title: "Conversation 1",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 1634000000,
								Author:     models.Author{Role: "user"},
								Content:    models.Content{Parts: []interface{}{"Hello!"}},
							},
						},
						"key2": {
							Message: models.Message{
								CreateTime: 1634000100,
								Author:     models.Author{Role: "assistant"},
								Content:    models.Content{Parts: []interface{}{"Hi there!"}},
							},
						},
					},
				},
			},
			titles: []string{"Conversation 1"},
			expected: []models.Snippet{
				{CreateTime: convertTime(1634000000), Role: "user", Content: "Hello!"},
				{CreateTime: convertTime(1634000100), Role: "assistant", Content: "Hi there!"},
			},
		},
		{
			name: "ignore invalid role",
			input: []models.Conversation{
				{
					Title: "Conversation 2",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 1634000000,
								Author:     models.Author{Role: "system"}, // invalid role
								Content:    models.Content{Parts: []interface{}{"System message"}},
							},
						},
					},
				},
			},
			titles:   []string{"Conversation 2"},
			expected: []models.Snippet{},
		},
		{
			name: "ignore empty content",
			input: []models.Conversation{
				{
					Title: "Conversation 2",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 1634000100,
								Author:     models.Author{Role: "user"},
								Content:    models.Content{Parts: []interface{}{}}, // empty content
							},
						},
					},
				},
			},
			titles:   []string{"Conversation 2"},
			expected: []models.Snippet{},
		},
		{
			name: "ignore messages with zero timestamp",
			input: []models.Conversation{
				{
					Title: "Conversation 3",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 0, // zero timestamp
								Author:     models.Author{Role: "user"},
								Content:    models.Content{Parts: []interface{}{"Zero timestamp"}},
							},
						},
						"key2": {
							Message: models.Message{
								CreateTime: 1634000200,
								Author:     models.Author{Role: "assistant"},
								Content:    models.Content{Parts: []interface{}{"Valid timestamp"}},
							},
						},
					},
				},
			},
			titles: []string{"Conversation 3"},
			expected: []models.Snippet{
				{CreateTime: convertTime(1634000200), Role: "assistant", Content: "Valid timestamp"},
			},
		},
		{
			name: "ignore messages when content is not type string",
			input: []models.Conversation{
				{
					Title: "Conversation 3",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 1634000100, // zero timestamp
								Author:     models.Author{Role: "user"},
								Content:    models.Content{Parts: []interface{}{map[string]interface{}{"key_001": "metadata here"}}},
							},
						},
					},
				},
			},
			titles:   []string{"Conversation 3"},
			expected: []models.Snippet{},
		},
		{
			name: "test sorting by CreateTime",
			input: []models.Conversation{
				{
					Title: "Conversation 4",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 1634000100,
								Author:     models.Author{Role: "user"},
								Content:    models.Content{Parts: []interface{}{"Second"}},
							},
						},
						"key2": {
							Message: models.Message{
								CreateTime: 1634000000,
								Author:     models.Author{Role: "assistant"},
								Content:    models.Content{Parts: []interface{}{"First"}},
							},
						},
					},
				},
			},
			titles: []string{"Conversation 4"},
			expected: []models.Snippet{
				{CreateTime: convertTime(1634000000), Role: "assistant", Content: "First"},
				{CreateTime: convertTime(1634000100), Role: "user", Content: "Second"},
			},
		},
		{
			name: "ignore conversation",
			input: []models.Conversation{
				{
					Title: "Conversation 4",
					Mapping: map[string]models.Mapping{
						"key1": {
							Message: models.Message{
								CreateTime: 1634000100,
								Author:     models.Author{Role: "user"},
								Content:    models.Content{Parts: []interface{}{"Second"}},
							},
						},
						"key2": {
							Message: models.Message{
								CreateTime: 1634000000,
								Author:     models.Author{Role: "assistant"},
								Content:    models.Content{Parts: []interface{}{"First"}},
							},
						},
					},
				},
			},
			titles:   []string{},
			expected: []models.Snippet{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapper.MapToSnippets(tc.input, tc.titles)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapToCSVRow(t *testing.T) {
	testCases := []struct {
		name        string
		input       models.Snippet
		expectedRow []string
	}{
		{
			name: "valid input",
			input: models.Snippet{
				CreateTime: convertTime(1634000000),
				Role:       "user",
				Content:    "Hello there!",
			},
			expectedRow: []string{"2021-10-12 00:53:20 +0000 UTC", "user", "Hello there!"},
		},
		{
			name: "empty role and content",
			input: models.Snippet{
				CreateTime: convertTime(1634000100),
				Role:       "",
				Content:    "",
			},
			expectedRow: []string{"2021-10-12 00:55:00 +0000 UTC", "", ""},
		},
		{
			name: "negative timestamp",
			input: models.Snippet{
				CreateTime: convertTime(-12345.678),
				Role:       "assistant",
				Content:    "Negative timestamp",
			},
			expectedRow: []string{"1969-12-31 20:34:15 +0000 UTC", "assistant", "Negative timestamp"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			row := mapper.MapToCSVRow(tc.input)
			assert.Equal(t, tc.expectedRow, row)
		})
	}
}

func convertTime(timeUnix float64) string {
	return time.Unix(int64(timeUnix), 0).UTC().String()
}
