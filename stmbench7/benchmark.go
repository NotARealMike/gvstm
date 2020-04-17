package main

import (
    "fmt"
    "gvstm/gvstm"
    "gvstm/stm"
    gvstm2 "gvstm/stmbench7/impl/gvstm"
    "gvstm/stmbench7/interfaces"
)

type benchmark interface {

}

type benchmarkImpl struct {
    synchMethodInitialiser interfaces.SynchMethodInitialiser
    reexecution bool
    gvstm bool
}

func createBenchmark(args []string, reexecution bool) benchmark {
    b := &benchmarkImpl{
        reexecution:reexecution,
    }
    b.parseCommandLineParameters(args)
    if !reexecution {
        b.printRunTimeParametersInformation()
    }
    b.generateOperationCDF()
    b.setupStructures()
    return b
}

func (b *benchmarkImpl) parseCommandLineParameters(args []string) {
    b.gvstm = true
    b.synchMethodInitialiser = gvstm2.GVSTMInitialiser
    if !b.reexecution {
        interfaces.SetFactories(b.synchMethodInitialiser)
    }
}

func (b *benchmarkImpl) printRunTimeParametersInformation() {}
func (b *benchmarkImpl) generateOperationCDF() {}
func (b *benchmarkImpl) setupStructures() {
    fmt.Println("Setup start...")
    b.createSetup()
    fmt.Println("Setup completed.")
}

func (b *benchmarkImpl) createSetup() interfaces.Setup {
    if b.gvstm {
        var setup interfaces.Setup
        gvstm.Atomic(func(tx stm.Transaction) {
            setup = interfaces.NewSetup(tx)
        })
        return setup
    }
    return interfaces.NewSetup(nil)
}