package interfaces

import (
	. "github.com/NotARealMike/gvstm/stm"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"math/rand"
)

type CompositePartBuilder interface {
	CreateAndRegisterCompositePart(tx Transaction) (CompositePart, OpFailedError)
	UnregisterAndRecycleCompositePart(tx Transaction, compositePart CompositePart)
}

type compositePartBuilderImpl struct {
	idPool             IDPool
	compositePartIndex Index
	documentBuilder    documentBuilder
	atomicPartBuilder  atomicPartBuilder
}

func newCompositePartBuilder(tx Transaction, compositePartIndex, documentIndex, atomicPartIndex, buildDateIndex Index) CompositePartBuilder {
	return &compositePartBuilderImpl{
		idPool:             BEFactory.CreateIDPool(tx, internal.MaxCompParts),
		compositePartIndex: compositePartIndex,
		documentBuilder:    newDocumentBuilder(tx, documentIndex),
		atomicPartBuilder:  newAtomicPartBuilder(tx, atomicPartIndex, buildDateIndex),
	}
}

func (cpb *compositePartBuilderImpl) CreateAndRegisterCompositePart(tx Transaction) (CompositePart, OpFailedError) {
	id, err := cpb.idPool.GetID(tx)
	if err != nil {
		return nil, err
	}

	typ := createType()
	var date int
	if rand.Int() < internal.YoungCompositePartFraction {
		date = createBuildDate(internal.MinYoungCompositePartDate, internal.MaxYoungCompositePartDate)
	} else {
		date = createBuildDate(internal.MinOldCompositePartDate, internal.MaxOldCompositePartDate)
	}

	document, err := cpb.documentBuilder.createAndRegisterDocument(tx, id)
	if err != nil {
		return nil, err
	}

	parts := make([]AtomicPart, internal.NumAtomsPerCompPart)
	for i := range parts {
		part, err := cpb.atomicPartBuilder.createAndRegisterAtomicPart(tx)
		if err == nil {
			parts[i] = part
			continue
		}
		// If creating an atomic part caused an error then we must unregister the
		// document and all atomic parts previously created.
		cpb.documentBuilder.unregisterAndRecycleDocument(tx, document)
		for j := range parts {
			if parts[j] != nil {
				cpb.atomicPartBuilder.unregisterAndRecycleAtomicPart(tx, parts[j])
			}
		}
		cpb.idPool.PutUnusedID(tx, id)
		return nil, err
	}

	cpb.createConnections(tx, parts)
	compositePart := DOFactory.CreateCompositePart(tx, id, typ, date, document)
	for _, part := range parts {
		compositePart.AddPart(tx, part)
	}
	cpb.compositePartIndex.Put(tx, id, compositePart)
	return compositePart, nil
}

func (cpb *compositePartBuilderImpl) UnregisterAndRecycleCompositePart(tx Transaction, compositePart CompositePart) {
	id := compositePart.GetId(tx)
	cpb.compositePartIndex.Remove(tx, id)

	cpb.documentBuilder.unregisterAndRecycleDocument(tx, compositePart.GetDocumentation(tx))

	for _, part := range compositePart.GetParts(tx).ToSlice(tx) {
		cpb.atomicPartBuilder.unregisterAndRecycleAtomicPart(tx, part.(AtomicPart))
	}

	for _, owner := range compositePart.GetUsedIn(tx).ToSlice() {
		// A composite part may appear multiple times in a base assembly
		for owner.(BaseAssembly).RemoveComponent(tx, compositePart) {
		}
	}

	compositePart.ClearPointers(tx)
	cpb.idPool.PutUnusedID(tx, id)
}

func (cpb *compositePartBuilderImpl) createConnections(tx Transaction, parts []AtomicPart) {
	// First, make all atomic parts be connected in a ring
	// (so that the resulting graph is fully connected)
	for i := 0; i < internal.NumAtomsPerCompPart; i++ {
		dest := (i + 1) % internal.NumAtomsPerCompPart
		parts[i].ConnectTo(tx, parts[dest], createType(), rand.Intn(internal.XYRange)+1)
	}

	// Then add other connections randomly, taking into account
	// the NumConnPerAtomic parameter. The procedure is non-deterministic
	// but it should eventually terminate.
	for i := 0; i < internal.NumAtomsPerCompPart; i++ {
		part := parts[i]
		for part.GetNumToConnections(tx) < internal.NumConnectionsPerAtomicPart {
			dest := rand.Intn(internal.NumAtomsPerCompPart)
			parts[i].ConnectTo(tx, parts[dest], createType(), rand.Intn(internal.XYRange)+1)
		}
	}
}
