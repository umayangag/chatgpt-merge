package writer

import (
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"encoding/csv"
	"os"
)

func ToCSV(filename string, snippets []models.Snippet) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Timestamp", "Role", "Content"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, snippet := range snippets {
		row, err := mapper.ToCSVRow(snippet)
		if err != nil {
			return err
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
