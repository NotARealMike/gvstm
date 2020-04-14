package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
)

type AssemblyBuilder struct {
    baseAssemblyIDPool, complexAssemblyIDPool IDPool
    baseAssemblyIndex, complexAssemblyIndex Index
}

func NewAssemblyBuilder(tx Transaction, baseAssemblyIndex, complexAssemblyIndex Index) *AssemblyBuilder {
    return &AssemblyBuilder{
        baseAssemblyIDPool:    beFactory.CreateIDPool(tx, internal.MaxBaseAssemblies),
        complexAssemblyIDPool: beFactory.CreateIDPool(tx, internal.MaxComplexAssemblies),
        baseAssemblyIndex:     baseAssemblyIndex,
        complexAssemblyIndex:  complexAssemblyIndex,
    }
}

func (ab *AssemblyBuilder) CreateAndRegisterAssembly(tx Transaction, module Module, superAssembly ComplexAssembly) (Assembly, *OpFailedError) {
    if superAssembly == nil || superAssembly.GetLevel(tx) > 2 {
        return ab.createAndRegisterComplexAssembly(tx, module, superAssembly)
    }
    return ab.createAndRegisterBaseAssembly(tx, module, superAssembly)
}

func (ab *AssemblyBuilder) UnregisterAndRecycleAssembly(tx Transaction, assembly Assembly) {
    switch t := assembly.(type) {
    case BaseAssembly :
        ab.unregisterAndRecycleBaseAssembly(tx, t)
    case ComplexAssembly :
        ab.unregisterAndRecycleComplexAssembly(tx, t)
    }
}

func (ab *AssemblyBuilder) createAndRegisterBaseAssembly(tx Transaction, module Module, superAssembly ComplexAssembly) (BaseAssembly, *OpFailedError) {
    id, err := ab.baseAssemblyIDPool.GetID(tx)
    if err != nil {
        return nil, err
    }

    date := createBuildDate(internal.MinAssDate, internal.MaxAssDate)
    baseAssembly := doFactory.CreateBaseAssembly(tx, id, createType(), date, module, superAssembly)

    ab.baseAssemblyIndex.Put(tx, id, baseAssembly)
    superAssembly.AddSubAssembly(tx, baseAssembly)
    return baseAssembly, nil
}

func (ab *AssemblyBuilder) createAndRegisterComplexAssembly(tx Transaction, module Module, superAssembly ComplexAssembly) (ComplexAssembly, *OpFailedError) {
    id, err := ab.complexAssemblyIDPool.GetID(tx)
    if err != nil {
        return nil, err
    }

    date := createBuildDate(internal.MinAssDate, internal.MaxAssDate)
    complexAssembly := doFactory.CreateComplexAssembly(tx, id, createType(), date, module, superAssembly)

    for i := 0; i < internal.NumSubAssemblies; i++ {
        _, err := ab.CreateAndRegisterAssembly(tx, module, complexAssembly)
        if err != nil {
            for _, subAssembly := range complexAssembly.GetSubAssemblies(tx).ToSlice() {
                ab.UnregisterAndRecycleAssembly(tx, subAssembly.(Assembly))
            }
            ab.complexAssemblyIDPool.PutUnusedID(tx, id)
            complexAssembly.ClearPointers(tx)
            return nil, err
        }
    }

    ab.complexAssemblyIndex.Put(tx, id, complexAssembly)
    if superAssembly != nil {
        superAssembly.AddSubAssembly(tx, complexAssembly)
    }
    return complexAssembly, nil
}

func (ab *AssemblyBuilder) unregisterAndRecycleBaseAssembly(tx Transaction, assembly BaseAssembly) {
    assemblyID := assembly.GetId(tx)
    ab.baseAssemblyIndex.Remove(tx, assemblyID)

    assembly.GetSuperAssembly(tx).RemoveSubAssembly(tx, assembly)

    components := assembly.GetComponents(tx).Clone()
    for _, component := range components.ToSlice() {
        assembly.RemoveComponent(tx, component.(CompositePart))
    }

    assembly.ClearPointers(tx)
    ab.baseAssemblyIDPool.PutUnusedID(tx, assembly.GetId(tx))
}

func (ab *AssemblyBuilder) unregisterAndRecycleComplexAssembly(tx Transaction, assembly ComplexAssembly) {
    currentLevel := assembly.GetLevel(tx)
    superAssembly := assembly.GetSuperAssembly(tx)

    // TODO: attempting to remove the root assembly is a bug according to the application logic. Should panic.

    superAssembly.RemoveSubAssembly(tx, assembly)

    subAssemblies := assembly.GetSubAssemblies(tx)
    if currentLevel > 2 {
        for _, subAssembly := range subAssemblies.ToSlice() {
            ab.unregisterAndRecycleComplexAssembly(tx, subAssembly.(ComplexAssembly))
        }
    } else {
        for _, subAssembly := range subAssemblies.ToSlice() {
            ab.unregisterAndRecycleBaseAssembly(tx, subAssembly.(BaseAssembly))
        }
    }

    id := assembly.GetId(tx)
    ab.complexAssemblyIndex.Remove(tx, id)

    assembly.ClearPointers(tx)
    ab.complexAssemblyIDPool.PutUnusedID(tx, id)

}