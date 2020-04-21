package main

import (
    "fmt"
    . "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
    "gvstm/stmbench7/operations"
    "math"
    "math/rand"
    "os"
    "sync"
    "sync/atomic"
    "time"
)

type benchThread interface {
    run(wg *sync.WaitGroup)
    stopThread()
    opsTTC(opIDIndex, ttc int) int
    opsHighTTCLog(opIDIndex, ttc int) int
    successful(opIDIndex int) int
    failed(opIDIndex int) int
    maxTTC(opIDIndex int) int
}

type benchThreadImpl struct {
    threadNum int
    stop atomic.Value
    operationCDF []float64
    executors []operations.OperationExecutor
    successfulOperations, failedOperations []int
    operationsTTC, operationsHighTTCLog    [][]int
}

func createBenchThread(setup Setup, operationCDF []float64, threadNum int) benchThread {
    numOfOperations := len(operations.OperationIDs)
    b := &benchThreadImpl{
        threadNum:            threadNum,
        stop:                 atomic.Value{},
        operationCDF:         operationCDF,
        executors:            make([]operations.OperationExecutor, numOfOperations),
        successfulOperations: make([]int, numOfOperations),
        failedOperations:     make([]int, numOfOperations),
        operationsTTC:        make([][]int, numOfOperations),
        operationsHighTTCLog: make([][]int, numOfOperations),
    }
    b.stop.Store(false)
    for i := range b.operationsTTC {
        b.operationsTTC[i] = make([]int, internal.MaxLowTTC + 1)
        b.operationsHighTTCLog[i] = make([]int, internal.HighTTCEntries)
    }
    b.createOperations(setup)
    return b
}

func (bt *benchThreadImpl) run(wg *sync.WaitGroup) {
    defer wg.Done()
    for bt.shouldContinue() {
        operationNumber := bt.getNextOperationNumber()
        executor := bt.executors[operationNumber]

        startTime := time.Now()
        _, err := executor.Execute()
        endTime := time.Now()

        if err == nil {
            bt.successfulOperations[operationNumber]++
            ttc := int(endTime.Sub(startTime).Nanoseconds()) / 1000000
            if ttc <= internal.MaxLowTTC {
                bt.operationsTTC[operationNumber][ttc]++
            } else {
                logHighTTC := (math.Log(float64(ttc)) - math.Log(float64(internal.MaxLowTTC))) / math.Log(internal.HighTTCLogBase)
                intLogHighTTC := int(math.Min(logHighTTC, float64(internal.HighTTCEntries)))
                bt.operationsHighTTCLog[operationNumber][intLogHighTTC]++
            }
        } else {
            bt.failedOperations[operationNumber]++
        }
    }
    fmt.Fprintf(os.Stderr, "Thread #%d finished.\n", bt.threadNum)
}

func (bt *benchThreadImpl) shouldContinue() bool {
    return !bt.stop.Load().(bool)
}

func (bt *benchThreadImpl) stopThread() {
    bt.stop.Store(true)
}

func (bt *benchThreadImpl) createOperations(setup Setup) {
    for i := range operations.OperationIDs {
        operation := operations.OperationIDs[i].CreateOperation(setup)
        bt.executors[i] = operations.OEFactory.CreateExecutor(operation)
    }
}

func (bt *benchThreadImpl) getNextOperationNumber() int {
    whichOperation := rand.Float64()
    operationNumber := 0
    for whichOperation >= bt.operationCDF[operationNumber] {
        operationNumber++
    }
    return operationNumber
}

func (bt *benchThreadImpl) opsTTC(opIDIndex, ttc int) int {
    return bt.operationsTTC[opIDIndex][ttc]
}

func (bt *benchThreadImpl) opsHighTTCLog(opIDIndex, ttc int) int {
    return bt.operationsHighTTCLog[opIDIndex][ttc]
}

func (bt *benchThreadImpl) successful(opIDIndex int) int {
    return bt.successfulOperations[opIDIndex]
}

func (bt *benchThreadImpl) failed(opIDIndex int) int {
    return bt.failedOperations[opIDIndex]
}

func (bt *benchThreadImpl) maxTTC(opIDIndex int) int {
    for logTTCIndex := internal.HighTTCEntries - 1 ; logTTCIndex >= 0 ; logTTCIndex-- {
        if bt.operationsHighTTCLog[opIDIndex][logTTCIndex] > 0 {
            return logTTCIndexToTTC(logTTCIndex)
        }
    }
    for ttc := internal.MaxLowTTC ; ttc >= 0 ; ttc-- {
        if bt.operationsTTC[opIDIndex][ttc] > 0 {
            return ttc
        }
    }
    return 0
}