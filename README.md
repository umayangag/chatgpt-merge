# ChatGPT-Merge

ChatGPT-Merge is a small Go CLI that converts exported ChatGPT conversations (JSON) into a clean CSV for analysis or sharing. It supports a dry-run mode to list conversation titles and selective merging using an include list.

## Quick Start

Prerequisites: Go 1.25+

```bash
# Clone and build
git clone https://github.com/umayangag/chatgpt-merge.git
cd chatgpt-merge
make build   # or: go build -o chatgpt-merge ./cmd

# Show version
./chatgpt-merge -version
```

## Usage

The CLI has two primary modes:

- Dry-run: outputs the list of conversation titles (no CSV is created)
- Merge: creates a CSV for selected conversations

Flags:
- `-dry`         Print conversation titles to the given output file (and stdout) without merging
- `-include`     Path to a text file containing the conversation titles to include (one per line)
- `-version`     Print the tool version and exit
- `-no-header`   Omit the CSV header row (default: header included)
- `-bom`         Write a UTF-8 BOM at the start of the CSV (useful for Excel)

Positional arguments:
- `<source.json>` Path to exported ChatGPT conversations JSON
- `<output>`      Path to output file (titles list in dry-run, CSV in merge)

### Examples

Dry-run (list titles to a file and stdout):

```bash
./chatgpt-merge -dry input/conversations.json output/titles.txt
```

Merge selected conversations (using titles from a file):

```bash
./chatgpt-merge -include input/include.txt input/conversations.json output/output.csv
```

Notes:
- The include file is line-based; empty/whitespace-only lines are ignored.
- Titles are matched exactly after trimming.

## CSV Semantics

- Columns: `Timestamp`, `Role`, `Content`
- Timestamp format: RFC3339 in UTC (e.g., `2021-10-12T00:53:20Z`)
- Ordering: rows are sorted chronologically by timestamp
- Defaults: header row is included; UTF‑8 BOM is not written

In a later phase we may expose flags for writer options (`-no-header`, `-bom`). Defaults today are header on, BOM off.

## Development

Common tasks are available via Makefile:

```bash
make fmt   # gofmt -s -w .
make lint  # go vet ./...
make test  # go test ./...
make build # go build -o chatgpt-merge ./cmd
make cover # coverage summary
```

Verification:
```bash
./chatgpt-merge -version
./chatgpt-merge -dry input/conversations.json output/titles.txt
```

## Contributing

- Create a feature branch (do not commit to main): `git checkout -b feat/your-change`
- Use Conventional Commits (e.g., `feat:`, `fix:`, `docs:`)
- Add/adjust tests where applicable
- Open a focused pull request

## License

Apache-2.0. See `LICENSE`.

