package operations

import (
    "gvstm/gvstm"
    . "gvstm/stm"
)

type gvstmOperationExecutor struct {
    operation Operation
}

func CreateGVSTMOperationExecutor(operation Operation) OperationExecutor {
    return &gvstmOperationExecutor{operation}
}

func (e *gvstmOperationExecutor) Execute() (result int, err error) {
    gvstm.Atomic(func(tx Transaction) {
        result, err = e.operation.Perform(tx)
    })
    return
}