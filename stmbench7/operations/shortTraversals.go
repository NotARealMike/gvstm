package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)


func createShortTraversal(id OperationID) Operation {
    return &operationImpl{
        id: id,
        f: func(tx Transaction) (i int, err error) {
            return 0, nil
        },
    }
}

func createShortTraversal1(setup Setup) Operation {
    return createShortTraversal(ST1)
}
func createShortTraversal2(setup Setup) Operation {
    return createShortTraversal(ST2)
}
func createShortTraversal3(setup Setup) Operation {
    return createShortTraversal(ST3)
}
func createShortTraversal4(setup Setup) Operation {
    return createShortTraversal(ST4)
}
func createShortTraversal5(setup Setup) Operation {
    return createShortTraversal(ST5)
}
func createShortTraversal6(setup Setup) Operation {
    return createShortTraversal(ST6)
}
func createShortTraversal7(setup Setup) Operation {
    return createShortTraversal(ST7)
}
func createShortTraversal8(setup Setup) Operation {
    return createShortTraversal(ST8)
}
func createShortTraversal9(setup Setup) Operation {
    return createShortTraversal(ST9)
}
func createShortTraversal10(setup Setup) Operation {
    return createShortTraversal(ST10)
}
