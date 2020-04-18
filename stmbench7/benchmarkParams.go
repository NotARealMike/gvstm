package main

import (
    "gvstm/stmbench7/interfaces"
)

type benchmarkParams struct {
    initialiser interfaces.SynchMethodInitialiser
    reexecution bool
    gvstm bool
}
