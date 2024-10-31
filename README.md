
# ChatGPT-Merge

ChatGPT-Merge is a command-line tool developed in Go that enables users to merge multiple ChatGPT conversation histories into a single, cohesive CSV file. This facilitates better context management, sharing, and analysis of ChatGPT interactions.

## Features

- **Merge Conversations**: Combine selected ChatGPT conversations into a unified CSV file.
- **Selective Inclusion**: Specify which conversations to include in the merge process.
- **Dry Run Mode**: Preview the list of conversations without performing the merge.

## Prerequisites

- **Go**: Ensure Go is installed on your system. You can download it from the [official Go website](https://go.dev/doc/install).

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/umayangag/chatgpt-merge.git
   ```

2. **Navigate to the Project Directory**:

   ```bash
   cd chatgpt-merge
   ```

3. **Build the Executable**:

   ```bash
   go mod tidy && go build -o chatgpt-merge cmd/main.go
   ```

   This command installs necessary dependencies and compiles the source code into an executable named `chatgpt-merge`.

## Usage

The tool offers two primary modes:

1. **Dry Run Mode**: Outputs a list of available conversations without merging.
2. **Merge Mode**: Merges selected conversations into a CSV file.

### Dry Run Mode

To preview and save the list of conversations:

```bash
./chatgpt-merge -dry path/to/conversations.json selected_titles.txt
```

- **Arguments**:
  - `-dry`: Activates dry run mode.
  - `path/to/conversations.json`: Path to the JSON file containing exported ChatGPT conversations.
  - `selected_titles.txt`: File where the list of conversation titles will be saved.

### Merge Mode

To merge selected conversations:

```bash
./chatgpt-merge -include=selected_titles.txt path/to/conversations.json output.csv
```

- **Arguments**:
  - `-include=selected_titles.txt`: Path to a text file containing titles of conversations to include in the merge.
  - `path/to/conversations.json`: Path to the JSON file containing exported ChatGPT conversations.
  - `output.csv`: Destination file for the merged CSV data.

**Note**: The `selected_titles.txt` file should list conversation titles, each on a new line, corresponding to the conversations you wish to merge.

## Workflow

1. **Export Conversations**: Export your ChatGPT conversations into a `conversations.json` file.
2. **Dry Run (Optional)**: Run the tool in dry run mode to obtain a list of conversation titles.
3. **Select Conversations**: Create a text file (`selected_titles.txt`) listing the titles of conversations you want to merge. Remove the titles of conversations you do not want to merge
4. **Merge Conversations**: Run the tool in merge mode to generate a CSV file containing the selected conversations.

## CSV Output Format

The resulting CSV file will have the following columns:

1. **Timestamp**: Unix timestamp of the message.
2. **Role**: Role of the message author (`assistant` or `user`).
3. **Content**: Text content of the message.

## Extending the Tool

Developers can extend the tool by:

- **Adding Filters**: Implement additional filtering criteria (e.g., by date range or keyword).
- **Supporting Other Formats**: Enable output in formats like JSON or XML.
- **Enhancing User Interface**: Develop a graphical user interface (GUI) for improved usability.
- **File Split Capabilities**: Split the output CSV file into multiple files if the size exceeds the upload filesize limits.

## Contributing

Contributions are welcome! To contribute:

1. **Fork the Repository**: Click the "Fork" button at the top right of the repository page.
2. **Create a New Branch**: Use `git checkout -b feature-branch-name`.
3. **Make Changes**: Implement your feature or fix.
4. **Commit Changes**: Use `git commit -m "Description of changes"`.
5. **Push to Branch**: Use `git push origin feature-branch-name`.
6. **Create a Pull Request**: Navigate to your forked repository and click the "New Pull Request" button.

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

Special thanks to the contributors and the open-source community for their support and feedback.

