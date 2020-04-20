package main

import (
    "gvstm/stmbench7/interfaces"
    "sync"
)

type benchThread interface {
    run(wg *sync.WaitGroup)
    shouldContinue() bool
    stopThread()
    createOperations(setup interfaces.Setup)
    getNextOperationNumber() int
}

type benchThreadImpl struct {

}

func createBenchThread(setup interfaces.Setup, operationCDF []float64, threadNum int) benchThread {
    return &benchThreadImpl{}
}

func (bt *benchThreadImpl) run(wg *sync.WaitGroup) {
    defer wg.Done()
}
func (bt *benchThreadImpl) shouldContinue() bool {return false}
func (bt *benchThreadImpl) stopThread() {}
func (bt *benchThreadImpl) createOperations(setup interfaces.Setup) {}
func (bt *benchThreadImpl) getNextOperationNumber() int {return 0}
