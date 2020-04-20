package main

import (
    "gvstm/stmbench7/interfaces"
    "time"
)

type benchmarkParams struct {
    initialiser interfaces.SynchMethodInitialiser
    reexecution bool
    gvstm bool

    numThreads int
    duration time.Duration

    readOnlyRatio int
    traversalsEnabled, structureModificationsEnabled, longReadWriteTraversalsEnabled bool
}
