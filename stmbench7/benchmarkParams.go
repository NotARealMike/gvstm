package main

import (
    "gvstm/stmbench7/impl/gvstm"
    "gvstm/stmbench7/impl/locks"
    "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
    "gvstm/stmbench7/operations"
    "time"
)

type benchmarkParams struct {
    initialiser interfaces.SynchMethodInitialiser
    executorFactory operations.OperationExecutorFactory
    reexecution bool
    gvstm bool
    syncType string

    numThreads int
    duration time.Duration

    readOnlyRatio int
    traversalsEnabled, structureModificationsEnabled, longReadWriteTraversalsEnabled bool
    printTTCHistograms bool
}

var (
    gvstmParamsPreset = &benchmarkParams{
        initialiser:                    gvstm.GVSTMInitialiser,
        executorFactory:                operations.OperationExecutorFactory{CreateExecutor: operations.CreateGVSTMOperationExecutor},
        reexecution:                    false,
        gvstm:                          true,
        syncType:                       "GVSTM",
        numThreads:                     4,
        duration:                       10 * time.Second,
        readOnlyRatio:                  internal.ReadDominatedWorkloadRORatio,
        traversalsEnabled:              true,
        structureModificationsEnabled:  false,
        longReadWriteTraversalsEnabled: true,
        printTTCHistograms:             false,
    }
    cgParamsPreset = &benchmarkParams{
        initialiser:                    locks.CGLocksInitialiser,
        executorFactory:                operations.OperationExecutorFactory{CreateExecutor: operations.CreateCGOperationExecutor},
        reexecution:                    false,
        gvstm:                          false,
        syncType:                       "Coarse grained locking",
        numThreads:                     4,
        duration:                       10 * time.Second,
        readOnlyRatio:                  internal.ReadDominatedWorkloadRORatio,
        traversalsEnabled:              true,
        structureModificationsEnabled:  false,
        longReadWriteTraversalsEnabled: true,
        printTTCHistograms:             false,
    }
    mgParamsPreset = &benchmarkParams{
        initialiser:                    locks.MGLocksInitialiser,
        executorFactory:                operations.OperationExecutorFactory{CreateExecutor: operations.CreateMGOperationExecutor},
        reexecution:                    false,
        gvstm:                          false,
        syncType:                       "Medium grained locking",
        numThreads:                     4,
        duration:                       10 * time.Second,
        readOnlyRatio:                  internal.ReadDominatedWorkloadRORatio,
        traversalsEnabled:              true,
        structureModificationsEnabled:  false,
        longReadWriteTraversalsEnabled: true,
        printTTCHistograms:             false,
    }
)
