package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)


func createTraversal(id OperationID) Operation {
    return &operationImpl{
        id: id,
        f: func(tx Transaction) (i int, err error) {
            return 0, nil
        },
    }
}

func createTraversal1(setup Setup) Operation {
    return createTraversal(T1)
}
func createTraversal2a(setup Setup) Operation {
    return createTraversal(T2a)
}
func createTraversal2b(setup Setup) Operation {
    return createTraversal(T2b)
}
func createTraversal2c(setup Setup) Operation {
    return createTraversal(T2c)
}
func createTraversal3a(setup Setup) Operation {
    return createTraversal(T3a)
}
func createTraversal3b(setup Setup) Operation {
    return createTraversal(T3b)
}
func createTraversal3c(setup Setup) Operation {
    return createTraversal(T3c)
}
func createTraversal4(setup Setup) Operation {
    return createTraversal(T4)
}
func createTraversal5(setup Setup) Operation {
    return createTraversal(T5)
}
func createTraversal6(setup Setup) Operation {
    return createTraversal(T6)
}
func createTraversalQ6(setup Setup) Operation {
    return createTraversal(Q6)
}
func createTraversalQ7(setup Setup) Operation {
    return createTraversal(Q7)
}
