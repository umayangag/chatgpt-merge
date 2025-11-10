package writer_test

import (
	"bytes"
	"chatgpt-merge/internal/models"
	"chatgpt-merge/internal/writer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteToCSV_DefaultOptions(t *testing.T) {
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
			expectedCSV: "Timestamp,Role,Content\n2021-10-12 00:53:20 +0000 UTC,user,Hello there!\n2021-10-12 00:55:00 +0000 UTC,assistant,How can I assist?\n",
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
			var buf bytes.Buffer
			err := writer.WriteToCSV(&buf, tc.mapToCSVRow, tc.snippets, writer.Options{IncludeHeader: true, WriteBOM: false})
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCSV, buf.String())
		})
	}
}

func TestWriteToCSV_NoHeader(t *testing.T) {
	var buf bytes.Buffer
	snippets := []models.Snippet{
		{CreateTime: convertTime(1634000000), Role: "user", Content: "Hello there!"},
	}
	mapRow := func(s models.Snippet) []string { return []string{s.CreateTime, s.Role, s.Content} }
	err := writer.WriteToCSV(&buf, mapRow, snippets, writer.Options{IncludeHeader: false, WriteBOM: false})
	assert.NoError(t, err)
	// No header, just the row
	assert.Equal(t, "2021-10-12 00:53:20 +0000 UTC,user,Hello there!\n", buf.String())
}

func TestWriteToCSV_WithBOM(t *testing.T) {
	var buf bytes.Buffer
	snippets := []models.Snippet{
		{CreateTime: convertTime(1634000000), Role: "user", Content: "Hello there!"},
	}
	mapRow := func(s models.Snippet) []string { return []string{s.CreateTime, s.Role, s.Content} }
	err := writer.WriteToCSV(&buf, mapRow, snippets, writer.Options{IncludeHeader: true, WriteBOM: true})
	assert.NoError(t, err)
	out := buf.Bytes()
	// Expect BOM at the beginning
	assert.GreaterOrEqual(t, len(out), 3)
	assert.Equal(t, []byte{0xEF, 0xBB, 0xBF}, out[:3])
	// Followed by header and one row
	rest := string(out[3:])
	expected := "Timestamp,Role,Content\n2021-10-12 00:53:20 +0000 UTC,user,Hello there!\n"
	assert.Equal(t, expected, rest)
}

func convertTime(timeUnix float64) string {
	return time.Unix(int64(timeUnix), 0).UTC().Format(time.RFC3339)
}
