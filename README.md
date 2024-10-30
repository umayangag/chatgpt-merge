
## Installation
1. Install golang [https://go.dev/doc/install]
2. Download the source code and run the following commands in the project directory

    `
    go mod tidy
    go build -o chatgpt-merge cmd/main.go
    `

## How to Use
1. Export the conversations from ChatGPT
2. Extract the files and get the path for `conversations.json` file
3. Run the following command with specifying input and output files to get the list of conversations

    `
    ./chatgpt-merge -dry conversations.json titles.txt
    `
4. Edit the `titles.txt` file to only include the conversations you want to merge
5. Run the following command with specifying input and output files to merge the selected conversations

    `
    ./chatgpt-merge -include=titles.txt conversations.json output.csv
    `