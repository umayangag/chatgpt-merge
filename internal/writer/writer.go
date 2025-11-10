package writer

import (
	"chatgpt-merge/internal/models"
	"encoding/csv"
	"io"
)

type CSVMapper func(s models.Snippet) []string

type Options struct {
	IncludeHeader bool
	WriteBOM      bool
}

func WriteToCSV(out io.Writer, mapToCSVRow CSVMapper, snippets []models.Snippet, opts Options) error {
	var w io.Writer = out
	if opts.WriteBOM {
		// Write UTF-8 BOM at the beginning of the stream
		if _, err := w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			return err
		}
	}

	csvWriter := csv.NewWriter(w)

	if opts.IncludeHeader {
		header := []string{"Timestamp", "Role", "Content"}
		if err := csvWriter.Write(header); err != nil {
			return err
		}
	}

	for _, snippet := range snippets {
		row := mapToCSVRow(snippet)
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
