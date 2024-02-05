package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
	"time"
)

type Result struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    *string   `json:"Test,omitempty"`
	Elapsed *float64  `json:"Elapsed,omitempty"`
	Output  *string   `json:"Output,omitempty"`
}

type TestName string
type TestResultsPassed map[TestName][]*Result
type TestResultsFailed map[TestName][]*Result
type TestResultsSkipped map[TestName][]*Result

func groupResultsByTestStatus(results []Result) (TestResultsPassed, TestResultsFailed, TestResultsSkipped) {
	failedResults := make(TestResultsFailed)
	passResults := make(TestResultsPassed)
	skipResults := make(TestResultsSkipped)

	var name TestName
	var outputs []*Result
	for _, d := range results {
		if d.Output != nil && stdoutFlag {
			fmt.Print(*d.Output)
		}
		if d.Test != nil {
			switch d.Action {
			case "run":
				outputs = []*Result{}
				name = TestName(*d.Test)
			case "output":
				dCopy := d
				outputs = append(outputs, &dCopy)
			case "fail":
				failedResults[name] = outputs
			case "skip":
				skipResults[name] = outputs
			case "pass":
				passResults[name] = outputs
			}
		}
	}

	return passResults, failedResults, skipResults
}

func printFailedResults(results TestResultsFailed) {
	separator := strings.Repeat("-", 120)
	color.Red(separator)
	color.Red("Failed")
	color.Red(separator)

	for testName, outputs := range results {
		for _, v := range outputs {
			if v.Output != nil {
				trimmed := strings.TrimSpace(*v.Output)
				if strings.Contains(trimmed, ".go") {
					color.Red("%s:%s %s\n", v.Package, testName, trimmed)
				}
			}
		}
	}
}

func printPassResults(result TestResultsPassed) {
	separator := strings.Repeat("-", 120)
	color.Green(separator)
	color.Green("Passed")
	color.Green(separator)

	for testName, outputs := range result {
		for _, v := range outputs {
			if v.Output != nil {
				trimmed := strings.TrimSpace(*v.Output)
				if strings.Contains(trimmed, "--- PASS") {
					color.Green("%s:%s %s\n", v.Package, testName, trimmed)
				}
			}
		}
	}
}

func printSkipResults(result TestResultsSkipped) {
	separator := strings.Repeat("-", 120)
	color.Blue(separator)
	color.Blue("Skipped")
	color.Blue(separator)

	for testName, outputs := range result {
		for _, v := range outputs {
			if v.Output != nil {
				trimmed := strings.TrimSpace(*v.Output)
				if strings.Contains(trimmed, "--- SKIP") {
					color.Blue("%s:%s %s\n", v.Package, testName, trimmed)
				}
			}
		}
	}
}

func loadTestResultsFromStdin() ([]Result, error) {
	scanner := bufio.NewScanner(os.Stdin)

	var results []Result
	for scanner.Scan() {
		var row Result
		if err := json.Unmarshal([]byte(scanner.Text()), &row); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}
		results = append(results, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func loadTestResultsFromFile(fileName string) ([]Result, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from file %s: %v", fileName, err)
	}

	var results []Result
	for scanner.Scan() {
		var row Result
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}
		results = append(results, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

var (
	commandLineFlagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fileName           string
	stdoutFlag         bool
	allFlag            bool
	ciFlag             bool
)

func run(args []string) int {
	commandLineFlagSet.StringVar(&fileName, "f", "", "Filename of the file containing the output from Go tests.")
	commandLineFlagSet.BoolVar(&stdoutFlag, "v", false, "Display the output of Go tests to stdout.")
	commandLineFlagSet.BoolVar(&allFlag, "a", false, "All (pass, fail, skip) results are output.")
	commandLineFlagSet.BoolVar(&ciFlag, "ci", false, "Set this flag to run in CI mode.")
	if err := commandLineFlagSet.Parse(args); err != nil {
		log.Fatal(err)
	}

	var results []Result
	var err error
	if fileName != "" {
		results, err = loadTestResultsFromFile(fileName)
		if err != nil {
			log.Fatal(err)
			return 1
		}
	} else {
		results, err = loadTestResultsFromStdin()
		if err != nil {
			log.Fatal(err)
			return 1
		}
	}

	if ciFlag {
		color.NoColor = false
	}

	passResult, failedResult, skipResult := groupResultsByTestStatus(results)
	printFailedResults(failedResult)
	if allFlag {
		printSkipResults(skipResult)
		printPassResults(passResult)
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
