package main

import (
	"github.com/NotARealMike/gvstm/stmbench7/interfaces"
	"github.com/NotARealMike/gvstm/stmbench7/operations"
	"time"
)

type benchmarkParams struct {
	initialiser     interfaces.SynchMethodInitialiser
	executorFactory operations.OperationExecutorFactory
	gvstm           bool
	syncType        string

	numThreads int
	duration   time.Duration

	readOnlyRatio                                                                    int
	traversalsEnabled, longReadWriteTraversalsEnabled, structureModificationsEnabled bool
	printTTCHistograms                                                               bool
}
