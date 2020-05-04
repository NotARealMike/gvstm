package main

import (
	"gvstm/stmbench7/interfaces"
	"gvstm/stmbench7/operations"
	"time"
)

type benchmarkParams struct {
	initialiser     interfaces.SynchMethodInitialiser
	executorFactory operations.OperationExecutorFactory
	reexecution     bool
	gvstm           bool
	syncType        string

	numThreads int
	duration   time.Duration

	readOnlyRatio                                                                    int
	traversalsEnabled, longReadWriteTraversalsEnabled, structureModificationsEnabled bool
	printTTCHistograms                                                               bool
}
