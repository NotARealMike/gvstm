package main

import (
    "fmt"
    "gvstm/gvstm"
    "gvstm/stm"
    "gvstm/stmbench7/correctness"
    "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
    "gvstm/stmbench7/operations"
    "os"
    "sync"
    "time"
)

type benchmark interface {
    createInitialClone()
    start()
    checkInvariants(initial bool) error
    checkOpacity() error
    showTTCHistograms()
    showStats()
}

type benchmarkImpl struct {
    params *benchmarkParams

    elapsedTime time.Duration
    operationCDF []float64
    benchThreads []benchThread
    setup interfaces.Setup
}

func createBenchmark(params *benchmarkParams) benchmark {
    b := &benchmarkImpl{params:params}
    if !b.params.reexecution {
        interfaces.SetFactories(b.params.initialiser)
        fmt.Fprintln(os.Stderr, header())
        s := runtimeParamsInfo(b.params)
        fmt.Fprintln(os.Stdout, s)
        fmt.Fprintln(os.Stderr, s)

    }
    b.generateOperationCDF()
    b.setupStructures()
    return b
}

func (b *benchmarkImpl) setupStructures() {
    fmt.Fprintln(os.Stderr, "Setup start...")
    if b.params.gvstm {
        gvstm.Atomic(func(tx stm.Transaction) {
            b.setup = interfaces.NewSetup(tx)
        })
    } else {
        b.setup = interfaces.NewSetup(nil)
    }
    b.benchThreads = make([]benchThread, b.params.numThreads)
    for i := 0 ; i < b.params.numThreads ; i++ {
        b.benchThreads[i] = createBenchThread(b.setup, b.operationCDF, i)
    }
    fmt.Fprintln(os.Stderr, "Setup completed.")
}

func (b *benchmarkImpl) createInitialClone() {}

func (b *benchmarkImpl) start() {
    fmt.Fprintln(os.Stderr, "\nBenchmark started.")
    // TODO: ThreadRandom
    startTime := time.Now()
    var wg sync.WaitGroup
    for _, t := range b.benchThreads {
        wg.Add(1)
        go t.run(&wg)
    }
    time.Sleep(b.params.duration)
    for _, t := range b.benchThreads {
        t.stopThread()
    }
    wg.Wait()
    b.elapsedTime = time.Now().Sub(startTime)
    fmt.Fprintln(os.Stderr, "Benchmark completed.\n")
}

func (b *benchmarkImpl) generateOperationCDF() {
    stRatio := float64(internal.ShortTraversalsRatio) / 100
    opRatio := float64(internal.OperationsRatio) / 100
    tvRatio, smRatio := 0., 0.
    if b.params.traversalsEnabled {
        tvRatio = float64(internal.TraversalsRatio) / 100
    }
    if b.params.structureModificationsEnabled {
        smRatio = float64(internal.StructuralModificationsRatio) / 100
    }

    readOnlyRatio := float64(b.params.readOnlyRatio) / 100
    updateRatio := 1 - readOnlyRatio

    roTraversalsRatio := tvRatio * readOnlyRatio
    rwTraversalsRatio := 0.
    if b.params.longReadWriteTraversalsEnabled {
        rwTraversalsRatio = tvRatio * updateRatio
    }

    sumRatios := roTraversalsRatio + rwTraversalsRatio + stRatio + opRatio + smRatio * updateRatio

    roTraversalsRatio /= sumRatios
    rwTraversalsRatio /= sumRatios
    stRatio /= sumRatios
    opRatio /= sumRatios
    smRatio /= sumRatios

    for _, opID := range operations.OperationIDs {
        opID.OpType.Count++
    }

    operations.TraversalRW.Probability = rwTraversalsRatio / float64(operations.TraversalRW.Count)
    operations.TraversalRO.Probability = roTraversalsRatio / float64(operations.TraversalRO.Count)
    operations.ShortTraversalRW.Probability = stRatio * updateRatio / float64(operations.ShortTraversalRW.Count)
    operations.ShortTraversalRO.Probability = stRatio * readOnlyRatio / float64(operations.ShortTraversalRO.Count)
    operations.OperationRW.Probability = opRatio * updateRatio / float64(operations.OperationRW.Count)
    operations.OperationRO.Probability = opRatio * readOnlyRatio / float64(operations.OperationRO.Count)
    operations.StructureModification.Probability = smRatio * updateRatio / float64(operations.StructureModification.Count)

    if !b.params.reexecution {
        fmt.Fprintln(os.Stderr, "Operation ratios [%]:")
        fmt.Fprintln(os.Stdout, "Operation ratios [%]:")
        for _, opType := range operations.OperationTypes {
            s := alignText(opType.Name, 23) + ": " + alignText(formatFloat(opType.Probability * float64(opType.Count) * 100), 6)
            fmt.Fprintln(os.Stderr, s)
            fmt.Fprintln(os.Stdout, s)
        }
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stdout)
    }

    operationProbabilities := make([]float64, len(operations.OperationIDs))
    for i := range operations.OperationIDs {
        operationProbabilities[i] = operations.OperationIDs[i].OpType.Probability
    }

    b.operationCDF = make([]float64, len(operations.OperationIDs))
    b.operationCDF[0] = operationProbabilities[0]
    for i := 1 ; i < len(operationProbabilities) ; i++ {
        b.operationCDF[i] = b.operationCDF[i-1] + operationProbabilities[i]
    }
    b.operationCDF[len(operationProbabilities)-1] = 1
}

func (b *benchmarkImpl) checkInvariants(initial bool) (err error) {
    if initial {
        fmt.Fprintln(os.Stdout, "Checking invariants (initial data structure): ")
        fmt.Fprintln(os.Stderr, "Checking invariants (initial data structure): ")

    } else {
        fmt.Fprintln(os.Stdout, "Checking invariants (final data structure): ")
        fmt.Fprintln(os.Stderr, "Checking invariants (final data structure): ")
    }
    if b.params.gvstm {
        gvstm.Atomic(func(tx stm.Transaction) {
            err = correctness.CheckInvariants(tx, b.setup, initial)
        })
    } else {
        err = correctness.CheckInvariants(nil, b.setup, initial)
    }
    if err == nil {
        fmt.Fprintln(os.Stderr, "Invariant check completed successfully!")
        fmt.Fprintln(os.Stdout, "Invariant check completed successfully!")
    }
    return
}

func (b *benchmarkImpl) checkOpacity() error {
    return nil
}

func (b *benchmarkImpl) showTTCHistograms() {}

func (b *benchmarkImpl) showStats() {}
