package main

import (
	"fmt"
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

	params := gvstmParamsPreset

	benchmark := createBenchmark(params)
	benchmark.createInitialClone()
	benchmark.start()
	if err := benchmark.checkInvariants(false); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	benchmark.checkOpacity()
	benchmark.showTTCHistograms()
	benchmark.showStats()
}
