# test-result CLI

The `test-result` CLI enhances the visibility and readability of Go test results by parsing JSON formatted output from the `go test` command. 
It provides a streamlined view of test outcomes, making it easier to understand test results at a glance. 
This tool is particularly useful for developers looking to quickly assess the status of their tests in continuous integration pipelines or local development environments.

## Features

- **JSON Parsing**: Directly parses JSON output from `go test -json ...`, providing a clear and structured display of test results.
- **Flexible Input**: Supports input from both standard input (stdin) and files, allowing for versatility in how test results are fed into the tool.
- **Customizable Output**: Offers options to display all test results or filter them based on their status (pass, fail, skip).
- **CI Integration**: Includes a `--ci` option for optimized output in Continuous Integration environments.

## Installation

TBD

## Usage
To use test-result, you can pipe the output of go test directly into it or specify a file containing the JSON output from previous go test executions.

From Standard Input (stdin)

```bash

go test -v -p 1 -json ./... | test-result

```
From a File

```bash
go test -v -p 1 -json ./... > results.json
test-result -f results.json
```

Command Line Options
```bash
-f <filename>: Specifies the filename containing the output from Go tests. If not provided, test-result will read from stdin.
-v: Displays the output of Go tests to stdout.
-a: All (pass, fail, skip) results are output.
--ci: Optimizes output for CI environments by automatically setting appropriate flags.
```

