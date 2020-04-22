package locks

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
)

type complexAssemblyImpl struct {
    Assembly
    subAssemblies smallSet
    level int
}

func newComplexAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) ComplexAssembly {
    lvl := internal.NumAssemblyLevels
    if superAssembly != nil {
        lvl = superAssembly.GetLevel(tx) - 1
    }
    return &complexAssemblyImpl{
        Assembly:      newAssemblyImpl(tx, id, typ, buildDate, module, superAssembly),
        subAssemblies: newSmallSetImpl(tx),
        level:         lvl,
    }
}

func (ca *complexAssemblyImpl) AddSubAssembly(tx Transaction, assembly Assembly) bool {
    // TODO: attempting to add a base assembly to a complex assembly whose level != 2 should panic.
    return ca.subAssemblies.Add(tx, assembly)
}

func (ca *complexAssemblyImpl) RemoveSubAssembly(tx Transaction, assembly Assembly) bool {
    return ca.subAssemblies.Remove(tx, assembly)
}

func (ca *complexAssemblyImpl) GetSubAssemblies(tx Transaction) ImmutableCollection {
    return NewImmutableCollectionImpl(ca.subAssemblies.ToSlice(tx))
}

func (ca *complexAssemblyImpl) GetLevel(tx Transaction) int {
    return ca.level
}

func (ca *complexAssemblyImpl) ClearPointers(tx Transaction) {
    ca.Assembly.ClearPointers(tx)
    ca.subAssemblies = nil
    ca.level = -1
}
