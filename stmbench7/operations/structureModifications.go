package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)


func createStructureModification(id OperationID) Operation {
    return &operationImpl{
        id: id,
        f: func(tx Transaction) (i int, err error) {
            return 0, nil
        },
    }
}

func createStructureModification1(setup Setup) Operation {
    return createStructureModification(SM1)
}
func createStructureModification2(setup Setup) Operation {
    return createStructureModification(SM2)
}
func createStructureModification3(setup Setup) Operation {
    return createStructureModification(SM3)
}
func createStructureModification4(setup Setup) Operation {
    return createStructureModification(SM4)
}
func createStructureModification5(setup Setup) Operation {
    return createStructureModification(SM5)
}
func createStructureModification6(setup Setup) Operation {
    return createStructureModification(SM6)
}
func createStructureModification7(setup Setup) Operation {
    return createStructureModification(SM7)
}
func createStructureModification8(setup Setup) Operation {
    return createStructureModification(SM8)
}
