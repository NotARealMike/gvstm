package main

import (
    "fmt"
    "strconv"
    "strings"
)

func header() string {
    return section("Unreleased version, there is no header.")
}

func runtimeParamsInfo(params *benchmarkParams) string {
    var b strings.Builder
    b.WriteString("Benchmark Parameters:\n" +
        fmt.Sprintf("  GVSTM: %t\n", params.gvstm) +
        fmt.Sprintf("  Reexecution: %t\n", params.reexecution),
    )
    return b.String()
}

func syntax() string {
    return "Unreleased version, there is no syntax.\n"
}

func section(title string) string {
    var b strings.Builder
    b.WriteString(line('-') + "\n")
    b.WriteString(title + "\n")
    b.WriteString(line('-') + "\n")
    return b.String()
}

func line(r rune) string {
    var b strings.Builder
    for i := 0 ; i < 79 ; i++ {
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
    for i := 0 ; i < padding ; i++ {
        b.WriteRune(' ')
    }
    b.WriteString(text)
    return b.String()
}

func formatFloat(f float64) string {
    return strconv.FormatFloat(f, 'f', 2, 64)
}