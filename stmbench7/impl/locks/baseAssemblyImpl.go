package locks

import (
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type baseAssemblyImpl struct {
	Assembly
	components bag
}

func newBaseAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) BaseAssembly {
	return &baseAssemblyImpl{
		Assembly:   newAssemblyImpl(tx, id, typ, buildDate, module, superAssembly),
		components: newBagImpl(tx),
	}
}

func (ba *baseAssemblyImpl) AddComponent(tx Transaction, component CompositePart) {
	ba.components.add(tx, component)
	component.AddAssembly(tx, ba)
}

func (ba *baseAssemblyImpl) RemoveComponent(tx Transaction, component CompositePart) bool {
	exists := ba.components.remove(tx, component)
	if !exists {
		return false
	}
	component.RemoveAssembly(tx, ba)
	return true
}

func (ba *baseAssemblyImpl) GetComponents(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(ba.components.toSlice(tx))
}

func (ba *baseAssemblyImpl) ClearPointer(tx Transaction) {
	ba.Assembly.ClearPointers(tx)
	ba.components = nil
}
