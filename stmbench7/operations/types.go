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
    Name                                           string
}

var (
    TraversalRW = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "TraversalRW",
    }
    TraversalRO = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "TraversalRO",
    }
    ShortTraversalRW = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "ShortTraversalRW",
    }
    ShortTraversalRO = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "ShortTraversalRO",
    }
    OperationRW = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "OperationRW",
    }
    OperationRO = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "OperationRO",
    }
    StructureModification = &OperationType{
        Count:                0,
        Probability:          0,
        SuccessfulOperations: 0,
        FailedOperations:     0,
        Maxttc:               0,
        Name:                 "StructureModification",
    }
)

var OperationTypes = []*OperationType{
    TraversalRW, TraversalRO, ShortTraversalRW, ShortTraversalRO, OperationRW, OperationRO, StructureModification,
}

type OperationID struct {
    OpType *OperationType
    CreateOperation func(setup Setup) Operation
    Name string
}
var (
    T1, T2a, T2b, T2c, T3a, T3b, T3c, T4, T5, T6, Q6, Q7,
    ST1, ST2, ST3, ST4, ST5, ST6, ST7, ST8, ST9, ST10,
    OP1, OP2, OP3, OP4, OP5, OP6, OP7, OP8, OP9, OP10, OP11, OP12, OP13, OP14, OP15,
    SM1, SM2, SM3, SM4, SM5, SM6, SM7, SM8 OperationID
)

var OperationIDs []OperationID

func init() {
    T1  = OperationID{TraversalRO, createTraversal1, "T1"}
    T2a = OperationID{TraversalRW, createTraversal2a, "T2a"}
    T2b = OperationID{TraversalRW, createTraversal2b, "T2b"}
    T2c = OperationID{TraversalRW, createTraversal2c, "T2c"}
    T3a = OperationID{TraversalRW, createTraversal3a, "T3a"}
    T3b = OperationID{TraversalRW, createTraversal3b, "T3b"}
    T3c = OperationID{TraversalRW, createTraversal3c, "T3c"}
    T4  = OperationID{TraversalRO, createTraversal4, "T4"}
    T5  = OperationID{TraversalRW, createTraversal5, "T5"}
    T6  = OperationID{TraversalRO, createTraversal6, "T6"}
    Q6  = OperationID{TraversalRO, createTraversalQ6, "Q6"}
    Q7  = OperationID{TraversalRO, createTraversalQ7, "Q7"}

    ST1  = OperationID{ShortTraversalRO, createShortTraversal1, "ST1"}
    ST2  = OperationID{ShortTraversalRO, createShortTraversal2, "ST2"}
    ST3  = OperationID{ShortTraversalRO, createShortTraversal3, "ST3"}
    ST4  = OperationID{ShortTraversalRO, createShortTraversal4, "ST4"}
    ST5  = OperationID{ShortTraversalRO, createShortTraversal5, "ST5"}
    ST6  = OperationID{ShortTraversalRW, createShortTraversal6, "ST6"}
    ST7  = OperationID{ShortTraversalRW, createShortTraversal7, "ST7"}
    ST8  = OperationID{ShortTraversalRW, createShortTraversal8, "ST8"}
    ST9  = OperationID{ShortTraversalRO, createShortTraversal9, "ST9"}
    ST10 = OperationID{ShortTraversalRW, createShortTraversal10, "ST10"}

    OP1  = OperationID{OperationRO, createOperation1, "OP1"}
    OP2  = OperationID{OperationRO, createOperation2, "OP2"}
    OP3  = OperationID{OperationRO, createOperation3, "OP3"}
    OP4  = OperationID{OperationRO, createOperation4, "OP4"}
    OP5  = OperationID{OperationRO, createOperation5, "OP5"}
    OP6  = OperationID{OperationRO, createOperation6, "OP6"}
    OP7  = OperationID{OperationRO, createOperation7, "OP7"}
    OP8  = OperationID{OperationRO, createOperation8, "OP8"}
    OP9  = OperationID{OperationRW, createOperation9, "OP9"}
    OP10 = OperationID{OperationRW, createOperation10, "OP10"}
    OP11 = OperationID{OperationRW, createOperation11, "OP11"}
    OP12 = OperationID{OperationRW, createOperation12, "OP12"}
    OP13 = OperationID{OperationRW, createOperation13, "OP13"}
    OP14 = OperationID{OperationRW, createOperation14, "OP14"}
    OP15 = OperationID{OperationRW, createOperation15, "OP15"}

    SM1 = OperationID{StructureModification, createStructureModification1, "SM1"}
    SM2 = OperationID{StructureModification, createStructureModification2, "SM2"}
    SM3 = OperationID{StructureModification, createStructureModification3, "SM3"}
    SM4 = OperationID{StructureModification, createStructureModification4, "SM4"}
    SM5 = OperationID{StructureModification, createStructureModification5, "SM5"}
    SM6 = OperationID{StructureModification, createStructureModification6, "SM6"}
    SM7 = OperationID{StructureModification, createStructureModification7, "SM7"}
    SM8 = OperationID{StructureModification, createStructureModification8, "SM8"}

    OperationIDs = []OperationID{
        T1, T2a, T2b, T2c, T3a, T3b, T3c, T4, T5, T6, Q6, Q7,
        ST1, ST2, ST3, ST4, ST5, ST6, ST7, ST8, ST9, ST10,
        OP1, OP2, OP3, OP4, OP5, OP6, OP7, OP8, OP9, OP10, OP11, OP12, OP13, OP14, OP15,
        SM1, SM2, SM3, SM4, SM5, SM6, SM7, SM8,
    }
}
