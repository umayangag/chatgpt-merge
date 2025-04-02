package writer_test

import (
	"chatgpt-merge/internal/models"
	"chatgpt-merge/internal/writer"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteToCSV(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		snippets    []models.Snippet
		mapToCSVRow writer.CSVMapper
		expectedCSV string
	}{
		{
			name: "successful write",
			snippets: []models.Snippet{
				{CreateTime: convertTime(1634000000), Role: "user", Content: "Hello there!"},
				{CreateTime: convertTime(1634000100), Role: "assistant", Content: "How can I assist?"},
			},
			mapToCSVRow: func(s models.Snippet) []string {
				return []string{s.CreateTime, s.Role, s.Content}
			},
			expectedCSV: "Timestamp,Role,Content\n1634000000,user,Hello there!\n1634000100,assistant,How can I assist?\n",
		},
		{
			name:     "empty snippets list",
			snippets: []models.Snippet{},
			mapToCSVRow: func(s models.Snippet) []string {
				return []string{s.CreateTime, s.Role, s.Content}
			},
			expectedCSV: "Timestamp,Role,Content\n", // only header row
		},
	}

	// Loop through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var csvOutput strings.Builder
			csvWriter := csv.NewWriter(&csvOutput)
			err := writer.WriteToCSV(csvWriter, tc.mapToCSVRow, tc.snippets)
			csvWriter.Flush()

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCSV, csvOutput.String())

		})
	}
}

func convertTime(timeUnix float64) string {
	return time.Unix(int64(timeUnix), 0).String()
}
