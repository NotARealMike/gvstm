package main

import (
    "fmt"
    "gvstm/stmbench7/impl/gvstm"
    "os"
    "time"
)

var (
    resultsDir = "/Users/ruicachopo/Desktop/stmbench7_results/"
    dateFormat = "2006-01-02 15:04"
)

func main() {
    t0 := time.Now()
    resultsFile, err := os.Create(resultsDir + t0.Format(dateFormat) + ".txt")
    if err != nil {
        panic(err)
    }
    defer resultsFile.Close()
    os.Stdout = resultsFile

    params := &benchmarkParams{
        initialiser: gvstm.GVSTMInitialiser,
        reexecution: false,
        gvstm:       true,
    }

    fmt.Fprintln(os.Stderr, header())

    benchmark := createBenchmark(params)
    benchmark.createInitialClone()
    benchmark.start()
    benchmark.checkInvariants(false)
    benchmark.checkOpacity()
    benchmark.showTTCHistograms()
    benchmark.showStats()
}