package main

import (
    "fmt"
    "os"
    "time"
)

var (
    resultsDir = "/Users/ruicachopo/Desktop/stmbench7_results/"
    dateFormat = "2006-01-02_15:04"
)

func main() {
    t0 := time.Now()
    resultsFile, err := os.Create(resultsDir + t0.Format(dateFormat) + ".txt")
    if err != nil {
        panic(err)
    }
    defer resultsFile.Close()
    resultsFile.WriteString("Running STMBench7 at " + t0.Format(dateFormat) + "\n")
    createBenchmark(nil, false)
    tEnd := time.Now()
    fmt.Println("Total time elapsed: " + tEnd.Sub(t0).String())
}