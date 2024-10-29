package mapper_test

import (
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToSnippets(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		input    []models.Conversation
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
			expected: []models.Snippet{
				{CreateTime: 1634000000, Role: "user", Content: "Hello!"},
				{CreateTime: 1634000100, Role: "assistant", Content: "Hi there!"},
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
			expected: []models.Snippet{
				{CreateTime: 1634000200, Role: "assistant", Content: "Valid timestamp"},
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
			expected: []models.Snippet{
				{CreateTime: 1634000000, Role: "assistant", Content: "First"},
				{CreateTime: 1634000100, Role: "user", Content: "Second"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapper.MapToSnippets(tc.input, []string{})
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
				CreateTime: 1634000000,
				Role:       "user",
				Content:    "Hello there!",
			},
			expectedRow: []string{"1634000000", "user", "Hello there!"},
		},
		{
			name: "empty role and content",
			input: models.Snippet{
				CreateTime: 1634000100,
				Role:       "",
				Content:    "",
			},
			expectedRow: []string{"1634000100", "", ""},
		},
		{
			name: "negative timestamp",
			input: models.Snippet{
				CreateTime: -12345.678,
				Role:       "assistant",
				Content:    "Negative timestamp",
			},
			expectedRow: []string{"-12345.678", "assistant", "Negative timestamp"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			row := mapper.MapToCSVRow(tc.input)
			assert.Equal(t, tc.expectedRow, row)
		})
	}
}
