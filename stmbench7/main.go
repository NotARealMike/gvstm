package main

import (
	"flag"
	"fmt"
	"github.com/NotARealMike/gvstm/stmbench7/impl/gvstm"
	"github.com/NotARealMike/gvstm/stmbench7/impl/locks"
	"github.com/NotARealMike/gvstm/stmbench7/operations"
	"os"
	"time"
)

var (
	dateFormat = "2006-01-02 15:04"
)

func main() {

	resultsDir := flag.String("outDir", "./", "output directory for benchmark results.")
	syncType := flag.String("sync", "cg", "synchronisation method to be benchmarked. Can be cg, mg or gvstm.")
	numThreads := flag.Int("threads", 4, "number of parallel threads.")
	duration := flag.Duration("duration", 10*time.Second, "duration of benchmark run.")
	readOnlyRatio := flag.Int("roRatio", 90, "percentage of read-only operations.")
	longTraversals := flag.Bool("traversals", true, "enable long traversals.")
	longRWTraversals := flag.Bool("rwTraversals", true, "enable long read-write traversals.")
	structureModifications := flag.Bool("sm", false, "enable structure modifications.")
	histograms := flag.Bool("hist", false, "print TTC histograms.")

	flag.Parse()

	params := &benchmarkParams{
		numThreads:                     *numThreads,
		duration:                       *duration,
		readOnlyRatio:                  *readOnlyRatio,
		traversalsEnabled:              *longTraversals,
		longReadWriteTraversalsEnabled: *longRWTraversals,
		structureModificationsEnabled:  *structureModifications,
		printTTCHistograms:             *histograms,
	}

	switch *syncType {
	case "cg":
		params.initialiser = locks.CGLocksInitialiser
		params.executorFactory = operations.OperationExecutorFactory{CreateExecutor: operations.CreateCGOperationExecutor}
		params.syncType = "Coarse grained locking"
	case "mg":
		params.initialiser = locks.MGLocksInitialiser
		params.executorFactory = operations.OperationExecutorFactory{CreateExecutor: operations.CreateMGOperationExecutor}
		params.syncType = "Medium grained locking"
	case "gvstm":
		params.initialiser = gvstm.GVSTMInitialiser
		params.executorFactory = operations.OperationExecutorFactory{CreateExecutor: operations.CreateGVSTMOperationExecutor}
		params.gvstm = true
		params.syncType = "gvstm"
	default:
		panic("invalid synchronisation method: " + *syncType)
	}

	t0 := time.Now()
	os.MkdirAll(*resultsDir, 0777)
	resultsFile, err := os.Create(*resultsDir + "/" + t0.Format(dateFormat) + ".txt")
	if err != nil {
		panic(err)
	}
	defer resultsFile.Close()
	os.Stdout = resultsFile

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
