package main

import (
	"bufio"
	"chatgpt-merge/internal/mapper"
	"chatgpt-merge/internal/models"
	"chatgpt-merge/internal/writer"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const version = "v0.1.0"

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		// Print error to stderr and exit with non-zero
		_, _ = fmt.Fprintln(os.Stderr, err) // best effort
		os.Exit(1)
	}
}

// run executes the CLI with provided args and I/O writers.
func run(argv []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("chatgpt-merge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	includeFrom := fs.String("include", "", "file path to list of conversation titles to include when merging (one per line)")
	isDryRun := fs.Bool("dry", false, "output the list of conversation titles without merging")
	showVersion := fs.Bool("version", false, "print version and exit")
	noHeader := fs.Bool("no-header", false, "omit the CSV header row")
	writeBOM := fs.Bool("bom", false, "write UTF-8 BOM at the start of the CSV (for Excel compatibility)")

	if err := fs.Parse(argv); err != nil {
		return err
	}

	if *showVersion {
		if _, err := fmt.Fprintf(stdout, "chatgpt-merge %s\n", version); err != nil {
			return err
		}
		return nil
	}

	args := fs.Args()
	// Validate common args
	if len(args) < 2 {
		return errors.New("invalid arguments: source path and output path are required. Usage: [-dry | -include <file>] <source.json> <output>")
	}
	source := args[0]
	output := args[1]

	// Read Source File
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("error reading source file: %w", err)
	}

	var conversations []models.Conversation
	if err := json.Unmarshal(data, &conversations); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if *isDryRun { // Output the list of conversations
		if err := dumpConversationList(conversations, output, stdout); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
		if _, err := fmt.Fprintln(stdout, "Conversation list successfully written to:", output); err != nil {
			return err
		}
		return nil
	}

	// Merge conversations
	if *includeFrom == "" {
		return errors.New("path to list of conversations is required: use -include or use -dry to output the list of conversations")
	}

	titlesData, err := os.ReadFile(*includeFrom)
	if err != nil {
		return fmt.Errorf("error reading include titles file: %w", err)
	}

	includeTitles := strings.Split(string(titlesData), "\n")
	snippets := mapper.MapToSnippets(conversations, includeTitles)

	// Create output CSV
	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	// Write CSV with options from flags (defaults: header on, BOM off)
	opts := writer.Options{IncludeHeader: !*noHeader, WriteBOM: *writeBOM}
	if err := writer.WriteToCSV(file, mapper.MapToCSVRow, snippets, opts); err != nil {
		_ = file.Close() // best-effort close before returning error
		return fmt.Errorf("error writing csv file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing output file: %w", err)
	}

	if _, err := fmt.Fprintln(stdout, "Data successfully extracted to:", output); err != nil {
		return err
	}
	return nil
}

func dumpConversationList(conversations []models.Conversation, output string, stdout io.Writer) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}

	writeErr := func() error {
		buf := bufio.NewWriter(file)
		for _, conversation := range conversations {
			if _, err := buf.WriteString(conversation.Title + "\n"); err != nil {
				return err
			}
			// mirror to stdout as well
			if _, err := fmt.Fprintln(stdout, conversation.Title); err != nil {
				return err
			}
		}
		return buf.Flush()
	}()

	closeErr := file.Close()

	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
