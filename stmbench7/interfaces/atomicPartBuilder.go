package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
    "math/rand"
)

type AtomicPartBuilder struct {
    idPool IDPool
    partIndex Index
    buildDateIndex Index
}

func NewAtomicPartBuilder(tx Transaction, partIndex, buildDateIndex Index) *AtomicPartBuilder {
    return &AtomicPartBuilder{
        idPool:         beFactory.CreateIDPool(tx, internal.MaxAtomicParts),
        partIndex:      partIndex,
        buildDateIndex: buildDateIndex,
    }
}

func (apb *AtomicPartBuilder) CreateAndRegisterAtomicPart(tx Transaction) (AtomicPart, *OpFailedError) {
    id, err := apb.idPool.GetID(tx)
    if err != nil {
        return nil, err
    }
    typ := createType()
    date := createBuildDate(internal.MinAtomicPartDate, internal.MaxAtomicPartDate)

    // TODO: Java version uses ThreadRandom. How random does this need to be.
    x := rand.Int()
    y := x + 1

    part := doFactory.CreateAtomicPart(tx, id, typ, date, x, y)

    apb.partIndex.Put(tx, id, part)
    AddAtomicPartToBuildDateIndex(tx, apb.buildDateIndex, part)
    return part, nil
}

func (apb *AtomicPartBuilder) UnregisterAndRecycleAtomicPart(tx Transaction, part AtomicPart) {
    id := part.GetId(tx)
    RemoveAtomicPartFromBuildDateIndex(tx, apb.buildDateIndex, part)
    apb.partIndex.Remove(tx, id)
    part.ClearPointers(tx)
    apb.idPool.PutUnusedID(tx, id)
}
