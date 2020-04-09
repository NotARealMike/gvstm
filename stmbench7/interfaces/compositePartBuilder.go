package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
    "math/rand"
)

type CompositePartBuilder struct {
    idPool IDPool
    compositePartIndex Index
    documentBuilder *DocumentBuilder
    atomicPartBuilder *AtomicPartBuilder
}

func NewCompositePartBuilder(tx Transaction, compositePartIndex, documentIndex, atomicPartIndex, buildDateIndex Index) *CompositePartBuilder {
    return &CompositePartBuilder{
        idPool:             beFactory.CreateIDPool(tx, internal.MaxCompParts),
        compositePartIndex: compositePartIndex,
        documentBuilder:    NewDocumentBuilder(tx, documentIndex),
        atomicPartBuilder:  NewAtomicPartBuilder(tx, atomicPartIndex, buildDateIndex),
    }
}

func (cpb *CompositePartBuilder) CreateAndRegisterCompositePart(tx Transaction) (CompositePart, *OpFailedError) {
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

    document, err := cpb.documentBuilder.CreateAndRegisterDocument(tx, id)
    if err != nil {
        return nil, err
    }

    parts := make([]AtomicPart, internal.NumAtomsPerCompPart)
    for i := range parts {
        part, err := cpb.atomicPartBuilder.CreateAndRegisterAtomicPart(tx)
        if err == nil {
            parts[i] = part
            continue
        }
        // If creating an atomic part caused an error then we must unregister the
        // document and all atomic parts previously created.
        cpb.documentBuilder.UnregisterAndRecycleDocument(tx, document)
        for j := range parts {
            if parts[j] != nil {
                cpb.atomicPartBuilder.UnregisterAndRecycleAtomicPart(tx, parts[j])
            }
        }
        cpb.idPool.PutUnusedID(tx, id)
        return nil, err
    }

    cpb.createConnections(tx, parts)
    compositePart := doFactory.CreateCompositePart(tx, id, typ, date, document)
    for _, part := range parts {
        compositePart.AddPart(tx, part)
    }
    cpb.compositePartIndex.Put(tx, id, compositePart)
    return compositePart, nil
}

func (cpb *CompositePartBuilder) UnregisterAndRecycleCompositePart(tx Transaction, compositePart CompositePart) {
    id := compositePart.GetId(tx)
    cpb.compositePartIndex.Remove(tx, id)

    cpb.documentBuilder.UnregisterAndRecycleDocument(tx, compositePart.GetDocumentation(tx))

    for _, part := range compositePart.GetParts(tx).ToSlice(tx) {
        cpb.atomicPartBuilder.UnregisterAndRecycleAtomicPart(tx, part.(AtomicPart))
    }

    for _, owner := range compositePart.GetUsedIn(tx).ToSlice() {
        // A composite part may appear multiple times in a base assembly
        for owner.(BaseAssembly).RemoveComponent(tx, compositePart) {}
    }

    compositePart.ClearPointers(tx)
    cpb.idPool.PutUnusedID(tx, id)
}

func (cpb *CompositePartBuilder) createConnections(tx Transaction, parts []AtomicPart) {
    // First, make all atomic parts be connected in a ring
    // (so that the resulting graph is fully connected)
    for i := 0; i < internal.NumAtomsPerCompPart; i++ {
        dest := (i + 1) % internal.NumAtomsPerCompPart
        // TODO: The java version uses ThreadRandom and a parameter for the RNG.
        parts[i].ConnectTo(tx, parts[dest], createType(), rand.Int())
    }

    // Then add other connections randomly, taking into account
    // the NumConnPerAtomic parameter. The procedure is non-deterministic
    // but it should eventually terminate.
    for i := 0; i < internal.NumAtomsPerCompPart; i++ {
        part := parts[i]
        for part.GetNumToConnections(tx) < internal.NumConnectionsPerAtomicPart {
            // TODO: The java version uses ThreadRandom and a parameter for the RNG.
            dest := rand.Int()
            parts[i].ConnectTo(tx, parts[dest], createType(), rand.Int())
        }
    }
}
