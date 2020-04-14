package gvstm

import (
    "gvstm/gvstm"
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type baseAssemblyImpl struct {
    Assembly
    components TVar
}

func newBaseAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) BaseAssembly {
    return &baseAssemblyImpl{
        Assembly:   newAssemblyImpl(tx, id, typ, buildDate, module, superAssembly),
        components: gvstm.CreateTVar(newBagImpl(tx)),
    }
}

func (ba *baseAssemblyImpl) AddComponent(tx Transaction, component CompositePart) {
    tx.Load(ba.components).(bagImpl).add(tx, component)
    component.AddAssembly(tx, ba)
}

func (ba *baseAssemblyImpl) RemoveComponent(tx Transaction, component CompositePart) bool {
    exists := tx.Load(ba.components).(bagImpl).remove(tx, component)
    if !exists {
        return false
    }
    component.RemoveAssembly(tx, ba)
    return true
}

func (ba *baseAssemblyImpl) GetComponents(tx Transaction) ImmutableCollection {
    return NewImmutableCollectionImpl(tx.Load(ba.components).(bagImpl).toSlice(tx))
}

func (ba *baseAssemblyImpl) ClearPointer(tx Transaction) {
    ba.Assembly.ClearPointers(tx)
    tx.Store(ba.components, nil)
}