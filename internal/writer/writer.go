package writer

import (
	"chatgpt-merge/internal/models"
	"encoding/csv"
)

type CSVMapper func(s models.Snippet) []string

func WriteToCSV(writer *csv.Writer, mapToCSVRow CSVMapper, snippets []models.Snippet) error {
	header := []string{"Timestamp", "Role", "Content"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, snippet := range snippets {
		row := mapToCSVRow(snippet)
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
