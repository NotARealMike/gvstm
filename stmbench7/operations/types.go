package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type Operation interface {
    ID() OperationID
    Perform(tx Transaction) (int, error)
}

type operationImpl struct {
    id OperationID
    f func(tx Transaction) (int, error)
}

func (o *operationImpl) ID() OperationID {
    return o.id
}

func (o *operationImpl) Perform(tx Transaction) (int, error) {
    return o.f(tx)
}

type OperationType struct {
    Count                                          int
    Probability                                    float64
    SuccessfulOperations, FailedOperations, Maxttc int
}

var (
    TraversalRW = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
    TraversalRO = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
    ShortTraversalRW = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
    ShortTraversalRO = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
    OperationRW = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
    OperationRO = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
    StructureModification = OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
    }
)

type OperationID struct {
    OpType OperationType
    CreateOperation func(setup Setup) Operation
}

var (
    T1  = OperationID{TraversalRO, createTraversal1}
    T2a = OperationID{TraversalRW, createTraversal2a}
    T2b = OperationID{TraversalRW, createTraversal2b}
    T2c = OperationID{TraversalRW, createTraversal2c}
    T3a = OperationID{TraversalRW, createTraversal3a}
    T3b = OperationID{TraversalRW, createTraversal3b}
    T3c = OperationID{TraversalRW, createTraversal3c}
    T4  = OperationID{TraversalRO, createTraversal4}
    T5  = OperationID{TraversalRW, createTraversal5}
    T6  = OperationID{TraversalRO, createTraversal6}
    Q6  = OperationID{TraversalRO, createTraversalQ6}
    Q7  = OperationID{TraversalRO, createTraversalQ7}

    ST1  = OperationID{ShortTraversalRO, createShortTraversal1}
    ST2  = OperationID{ShortTraversalRO, createShortTraversal2}
    ST3  = OperationID{ShortTraversalRO, createShortTraversal3}
    ST4  = OperationID{ShortTraversalRO, createShortTraversal4}
    ST5  = OperationID{ShortTraversalRO, createShortTraversal5}
    ST6  = OperationID{ShortTraversalRW, createShortTraversal6}
    ST7  = OperationID{ShortTraversalRW, createShortTraversal7}
    ST8  = OperationID{ShortTraversalRW, createShortTraversal8}
    ST9  = OperationID{ShortTraversalRO, createShortTraversal9}
    ST10 = OperationID{ShortTraversalRW, createShortTraversal10}

    OP1  = OperationID{OperationRO, createOperation1}
    OP2  = OperationID{OperationRO, createOperation2}
    OP3  = OperationID{OperationRO, createOperation3}
    OP4  = OperationID{OperationRO, createOperation4}
    OP5  = OperationID{OperationRO, createOperation5}
    OP6  = OperationID{OperationRO, createOperation6}
    OP7  = OperationID{OperationRO, createOperation7}
    OP8  = OperationID{OperationRO, createOperation8}
    OP9  = OperationID{OperationRW, createOperation9}
    OP10 = OperationID{OperationRW, createOperation10}
    OP11 = OperationID{OperationRW, createOperation11}
    OP12 = OperationID{OperationRW, createOperation12}
    OP13 = OperationID{OperationRW, createOperation13}
    OP14 = OperationID{OperationRW, createOperation14}
    OP15 = OperationID{OperationRW, createOperation15}

    SM1 = OperationID{StructureModification, createStructureModification1}
    SM2 = OperationID{StructureModification, createStructureModification2}
    SM3 = OperationID{StructureModification, createStructureModification3}
    SM4 = OperationID{StructureModification, createStructureModification4}
    SM5 = OperationID{StructureModification, createStructureModification5}
    SM6 = OperationID{StructureModification, createStructureModification6}
    SM7 = OperationID{StructureModification, createStructureModification7}
    SM8 = OperationID{StructureModification, createStructureModification8}
)

