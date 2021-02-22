package main

import (
	"bufio"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
)

// Definition of the function which has traceability
type function struct {
	docblock string
	filename string
	line int
	name string
	pkg string
	traces *[]trace
}

// Definition of the trace itself extracted from the doc block
type trace struct {
	category string
	description string
	epic string
}

func main() {
	fileSet := token.NewFileSet()

	pkgs, err := parser.ParseDir(
		fileSet,
		"./",
		func(info os.FileInfo) bool {
			// Only parse test files
			return strings.Contains(info.Name(), "_test.go")
		},
		parser.ParseComments,
	)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading comments from file: %v", err))
		os.Exit(1)
	}

	var funcs []function

	for pkg, file := range pkgs {
		documentation := doc.New(file, "./", 0)
		for _, definition := range documentation.Funcs {
			// Only test functions
			if !strings.HasPrefix(definition.Name, "Test") {
				continue
			}

			traces := extractTraces(definition.Doc)

			if len(*traces) == 0 {
				continue
			}

			funcs = append(funcs, function{
				docblock: definition.Doc,
				name:     definition.Name,
				filename: file.Name,
				pkg:      pkg,
				traces:   traces,
			})
		}
	}

	// Dump out functions struct for now
	for _, fnc := range funcs {
		fmt.Println("package: " + fnc.pkg)
		fmt.Println("function: " + fnc.name)
		fmt.Println("traces:")
		for _, trc := range *fnc.traces {
			fmt.Println(" - category: " + trc.category)
			fmt.Println("   epic: " + trc.epic)
			fmt.Println("   description: " + trc.description)
		}
		fmt.Println()
	}
}

// Capture the current trace string if it's not empty
func captureTrace(traces *[]trace, current string) {
	if current == "" {
		return
	}

	created := createTrace(current)

	if created != nil {
		*traces = append(*traces, *created)
	}
}

// Create a trace from a comment
func createTrace(comment string) *trace {
	// The pattern for traces, they will be in the godocs note format
	// e.g. TRACE(ABT-123): This is an example trace
	pattern := regexp.MustCompile(`([A-Z_]+)\((ABT-[0-9]+)\):?\s?(.*)`)
	matches := pattern.FindStringSubmatch(comment)

	if len(matches) > 0 && validateCategory(matches[1]) {
		return &trace{
			category: matches[1],
			description: matches[3],
			epic: matches[2],
		}
	}

	return nil
}

// Process each line of comments to capture traces
func extractTraces(comments string) *[]trace {
	scanner := bufio.NewScanner(strings.NewReader(comments))

	if scanner.Err() != nil {
		fmt.Println(fmt.Sprintf("Error encountered while trying to read comments: %v", scanner.Err()))
		os.Exit(1)
	}

	// Traces will be prefixed with a godoc note, e.g. FEATURE(ABT-111):
	traceStart := regexp.MustCompile(`^[A-Z_]+\(`)

	// Process the comment block and group the comments together with special
	// handling for multiline comments
	var current string
	traces := &[]trace{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// See if line is the start of a trace
		note := traceStart.FindString(line)
		if note != "" {
			// Capture the current trace if applicable, this is an indication
			// we've found the start of a new trace
			captureTrace(traces, current)

			// Record line start and loop
			current = line
			continue
		}

		// This is an empty line
		if line == "" {
			// Capture the current trace if applicable and reset current
			// to an empty string. Since this is a gap between comment lines
			// any comments found after this can't belong to the previous trace
			captureTrace(traces, current)
			current = ""

			// Loop
			continue
		}

		// This is a comment that isn't a note, it may be part of a
		// multiline comment. If there is a current trace going then
		// this line should be appended to the current trace, otherwise
		// it can safely be ignored
		if current != "" {
			current = fmt.Sprintf("%s %s", current, line)
		}
	}

	// If there is a current trace capture it so it's not missed
	captureTrace(traces, current)

	return traces
}

// Validate the category for a trace
func validateCategory(test string) bool {
	// Define valid trace types required for audit, traces starting with
	// any other category will be ignored
	categories := [2]string{
		"BUG",
		"FEATURE",
	}

	for _, category := range categories {
		if test == category {
			return true
		}
	}

	return false
}