package main

import (
    "fmt"
    "gvstm/gvstm"
    "gvstm/stm"
    "gvstm/stmbench7/correctness"
    "gvstm/stmbench7/interfaces"
    "os"
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

    setup interfaces.Setup
}

func createBenchmark(params *benchmarkParams) benchmark {
    b := &benchmarkImpl{params:params}
    if !b.params.reexecution {
        interfaces.SetFactories(b.params.initialiser)
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
    fmt.Fprintln(os.Stderr, "Setup completed.")
}

func (b *benchmarkImpl) createInitialClone() {}

func (b *benchmarkImpl) start() {}

func (b *benchmarkImpl) generateOperationCDF() {}

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
