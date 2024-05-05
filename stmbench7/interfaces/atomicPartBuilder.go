package interfaces

import (
	. "github.com/NotARealMike/gvstm/stm"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"math/rand"
)

type atomicPartBuilder interface {
	createAndRegisterAtomicPart(tx Transaction) (AtomicPart, OpFailedError)
	unregisterAndRecycleAtomicPart(tx Transaction, part AtomicPart)
}

type atomicPartBuilderImpl struct {
	idPool         IDPool
	partIndex      Index
	buildDateIndex Index
}

func newAtomicPartBuilder(tx Transaction, partIndex, buildDateIndex Index) atomicPartBuilder {
	return &atomicPartBuilderImpl{
		idPool:         BEFactory.CreateIDPool(tx, internal.MaxAtomicParts),
		partIndex:      partIndex,
		buildDateIndex: buildDateIndex,
	}
}

func (apb *atomicPartBuilderImpl) createAndRegisterAtomicPart(tx Transaction) (AtomicPart, OpFailedError) {
	id, err := apb.idPool.GetID(tx)
	if err != nil {
		return nil, err
	}
	typ := createType()
	date := createBuildDate(internal.MinAtomicPartDate, internal.MaxAtomicPartDate)

	x := rand.Int()
	y := x + 1

	part := DOFactory.CreateAtomicPart(tx, id, typ, date, x, y)

	apb.partIndex.Put(tx, id, part)
	AddAtomicPartToBuildDateIndex(tx, apb.buildDateIndex, part)
	return part, nil
}

func (apb *atomicPartBuilderImpl) unregisterAndRecycleAtomicPart(tx Transaction, part AtomicPart) {
	id := part.GetId(tx)
	RemoveAtomicPartFromBuildDateIndex(tx, apb.buildDateIndex, part)
	apb.partIndex.Remove(tx, id)
	part.ClearPointers(tx)
	apb.idPool.PutUnusedID(tx, id)
}
