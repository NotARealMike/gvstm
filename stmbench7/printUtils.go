package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func printOutAndErr(s string) {
	fmt.Fprint(os.Stdout, s)
	fmt.Fprint(os.Stderr, s)
}

func header() string {
	return section("Unreleased version, there is no header.")
}

func runtimeParamsInfo(params *benchmarkParams) string {
	var b strings.Builder
	b.WriteString(section("Benchmark parameters"))

	b.WriteString(fmt.Sprintf("Number of threads: %d\n", params.numThreads))
	b.WriteString(fmt.Sprintf("Length: %s\n", params.duration.String()))
	b.WriteString(fmt.Sprintf("Read percentage: %d\n", params.readOnlyRatio))
	b.WriteString(fmt.Sprintf("Synchronisation method: %s\n", params.syncType))
	b.WriteString(fmt.Sprintf("Long traversals enabled: %t\n", params.traversalsEnabled))
	b.WriteString(fmt.Sprintf("Long read-write traversals enabled: %t\n", params.longReadWriteTraversalsEnabled))
	b.WriteString(fmt.Sprintf("Structure modification enabled: %t\n\n", params.structureModificationsEnabled))

	return b.String()
}

func syntax() string {
	return "Unreleased version, there is no syntax.\n"
}

func section(title string) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(line('-') + "\n")
	b.WriteString(title + "\n")
	b.WriteString(line('-') + "\n\n")
	return b.String()
}

func line(r rune) string {
	var b strings.Builder
	for i := 0; i < 79; i++ {
		b.WriteRune(r)
	}
	return b.String()
}

func alignText(text string, width int) string {
	padding := width - len(text)
	if padding < 0 {
		return "ERROR: insufficient width!"
	}
	var b strings.Builder
	for i := 0; i < padding; i++ {
		b.WriteRune(' ')
	}
	b.WriteString(text)
	return b.String()
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
