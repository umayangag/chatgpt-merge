
## Installation
1. Install golang [https://go.dev/doc/install]
2. Download the source code and run the following command in the project directory

`
go mod tidy 
`

## How to Use
1. Export the conversations from ChatGPT
2. Extract the files and get the path for conversations.json file
3. Run the following command with specifying input and output files

Example:
`
go run cmd/main.go input/conversations.json output/output.csv
`