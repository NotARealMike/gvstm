package operations

import . "gvstm/stmbench7/interfaces"

type Operation interface {
    ID() OperationID
    Perform() (int, error)
}

type OperationType struct {
    count int
    probability float64
    successfulOperations, failedOperations, maxttc int
}

var (
    TraversalRW = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
    TraversalRO = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
    ShortTraversalRW = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
    ShortTraversalRO = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
    OperationRW = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
    OperationRO = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
    StructureModification = OperationType{
        count:                0,
        probability:          0,
        successfulOperations: 0,
        failedOperations:     0,
        maxttc:               0,
    }
)

type OperationID struct {
    typ OperationType
    creator func(setup Setup) Operation
}

var (
    T1 OperationID = OperationID{TraversalRO, createTraversal1}
    T2a OperationID = OperationID{TraversalRW, createTraversal2a}
    T2b OperationID = OperationID{TraversalRW, createTraversal2b}
    T2c OperationID = OperationID{TraversalRW, createTraversal2c}
    T3a OperationID = OperationID{TraversalRW, createTraversal3a}
    T3b OperationID = OperationID{TraversalRW, createTraversal3b}
    T3c OperationID = OperationID{TraversalRW, createTraversal3c}
    T4 OperationID = OperationID{TraversalRO, createTraversal4}
    T5 OperationID = OperationID{TraversalRW, createTraversal5}
    T6 OperationID = OperationID{TraversalRO, createTraversal6}
    Q6 OperationID = OperationID{TraversalRO, createTraversalQ6}
    Q7 OperationID = OperationID{TraversalRO, createTraversalQ7}

    ST1 OperationID = OperationID{ShortTraversalRO, createShortTraversal1}
    ST2 OperationID = OperationID{ShortTraversalRO, createShortTraversal2}
    ST3 OperationID = OperationID{ShortTraversalRO, createShortTraversal3}
    ST4 OperationID = OperationID{ShortTraversalRO, createShortTraversal4}
    ST5 OperationID = OperationID{ShortTraversalRO, createShortTraversal5}
    ST6 OperationID = OperationID{ShortTraversalRW, createShortTraversal6}
    ST7 OperationID = OperationID{ShortTraversalRW, createShortTraversal7}
    ST8 OperationID = OperationID{ShortTraversalRW, createShortTraversal8}
    ST9 OperationID = OperationID{ShortTraversalRO, createShortTraversal9}
    ST10 OperationID = OperationID{ShortTraversalRW, createShortTraversal10}

    OP1 OperationID = OperationID{OperationRO, createOperation1}
    OP2 OperationID = OperationID{OperationRO, createOperation2}
    OP3 OperationID = OperationID{OperationRO, createOperation3}
    OP4 OperationID = OperationID{OperationRO, createOperation4}
    OP5 OperationID = OperationID{OperationRO, createOperation5}
    OP6 OperationID = OperationID{OperationRO, createOperation6}
    OP7 OperationID = OperationID{OperationRO, createOperation7}
    OP8 OperationID = OperationID{OperationRO, createOperation8}
    OP9 OperationID = OperationID{OperationRW, createOperation9}
    OP10 OperationID = OperationID{OperationRW, createOperation10}
    OP11 OperationID = OperationID{OperationRW, createOperation11}
    OP12 OperationID = OperationID{OperationRW, createOperation12}
    OP13 OperationID = OperationID{OperationRW, createOperation13}
    OP14 OperationID = OperationID{OperationRW, createOperation14}
    OP15 OperationID = OperationID{OperationRW, createOperation15}

    SM1 OperationID = OperationID{StructureModification, createStructureModification1}
    SM2 OperationID = OperationID{StructureModification, createStructureModification2}
    SM3 OperationID = OperationID{StructureModification, createStructureModification3}
    SM4 OperationID = OperationID{StructureModification, createStructureModification4}
    SM5 OperationID = OperationID{StructureModification, createStructureModification5}
    SM6 OperationID = OperationID{StructureModification, createStructureModification6}
    SM7 OperationID = OperationID{StructureModification, createStructureModification7}
    SM8 OperationID = OperationID{StructureModification, createStructureModification8}
)

