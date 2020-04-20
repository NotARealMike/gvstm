package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
    "math/rand"
)


func createOperation(id OperationID) Operation {
    return &operationImpl{
        id: id,
        f: func(tx Transaction) (i int, err error) {
            return 0, nil
        },
    }
}

func createOperation1(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        count := 0
        for i := 0 ; i < 10 ; i++ {
            // TODO: ThreadRandom?
            partID := rand.Intn(internal.MaxAtomicParts) + 1
            part := setup.AtomicPartIDIndex().Get(tx, partID)
            if part == nil {
                continue
            }
            part.(AtomicPart).NullOperation(tx)
            count++
        }
        return count, nil
    }
    return &operationImpl{
        id: OP1,
        f:  f,
    }
}
func createOperation2(setup Setup) Operation {
    return createOperation(OP2)
}
func createOperation3(setup Setup) Operation {
    return createOperation(OP3)
}
func createOperation4(setup Setup) Operation {
    return createOperation(OP4)
}
func createOperation5(setup Setup) Operation {
    return createOperation(OP5)
}
func createOperation6(setup Setup) Operation {
    return createOperation(OP6)
}
func createOperation7(setup Setup) Operation {
    return createOperation(OP7)
}
func createOperation8(setup Setup) Operation {
    return createOperation(OP8)
}
func createOperation9(setup Setup) Operation {
    return createOperation(OP9)
}
func createOperation10(setup Setup) Operation {
    return createOperation(OP10)
}
func createOperation11(setup Setup) Operation {
    return createOperation(OP11)
}
func createOperation12(setup Setup) Operation {
    return createOperation(OP12)
}
func createOperation13(setup Setup) Operation {
    return createOperation(OP13)
}
func createOperation14(setup Setup) Operation {
    return createOperation(OP14)
}
func createOperation15(setup Setup) Operation {
    return createOperation(OP15)
}
