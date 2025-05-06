# Elastic Uploader

`Elastic Uploader` is designed to simplify the indexing of
structured text files into Elasticsearch. It supports creating the target index
with predefined settings and populating it by reading data from a local directory.

## Features

* **Elasticsearch Index Management:** Create a specific Elasticsearch index
with predefined mappings and settings.
* **Data Upload:** Read structured text files from a local directory and
upload their content in bulk to the target Elasticsearch index.

## Prerequisites

* Go (version 1.24.2 or later)
* A running Elasticsearch instance
* Required environment variables set

## Installation

1. Clone the repository:

   ```bash
   git clone [https://github.com/git-iac/elastic-uploader.git](https://github.com/git-iac/elastic-uploader.git)
   cd elastic-uploader
   ```

2. Tidy Go modules:

   ```bash
   make tidy
   ```

3. Build the application using the Makefile:

   ```bash
   make build 
   # or
   go run ./cmd/main.go # pass args here
   ```

   This will create the executable file at `./build/elastic-uploader`.

## Configuration

### Environment Variables

The application requires the following environment variables to connect to Elasticsearch:

* `ELASTIC_URL`: The URL of your Elasticsearch instance (e.g., `http://localhost:9200`).

* `ELASTIC_API_KEY`: An API key for authentication with Elasticsearch.

## Input Data

The application processes data by reading files from a specific directory within
the project structure.

* **Input Folder:** Create a folder named `sections` at the root level of the
project directory (`./sections`).
* **File Naming:** Files within the `sections` folder must follow the format
`NNNN_Section.txt`, where `NNNN` represents a zero-padded four-digit
number (e.g., `0001_Section.txt`, `0002_Section.txt`). The application will process files in numerical order based on this prefix.

## Index Settings

The configuration for the target Elasticsearch index, including its name,
mapping, and settings, is defined programmatically within the application's code.

* The index name and the JSON body (mapping, settings, etc.) for the index
creation request are located in the `pkg/constants.go` file. You may need to
modify this file to configure the index structure according to your requirements.

## Usage

Execute the compiled application binary from the root of your project directory.
The primary action the application performs is determined by the `-action` command-line flag.

```bash
./elastic-uploader [flags]
```

## Flags

* action (required): Specifies the main operation to perform.
  * create: Creates the Elasticsearch index using the settings defined in pkg/constants.go.
  * populate: Reads files from the ./sections directory, processes them, and
  uploads their content in bulk requests to the target Elasticsearch index.
* chunkSize (optional): Controls the number of file contents to process and
send in a single bulk request during the populate action.Defaults to 1000.
