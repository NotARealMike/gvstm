package gvstm

import (
	"github.com/NotARealMike/gvstm/gvstm"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
)

type complexAssemblyImpl struct {
	Assembly
	subAssemblies TVar
	level         TVar
}

func newComplexAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) ComplexAssembly {
	lvl := internal.NumAssemblyLevels
	if superAssembly != nil {
		lvl = superAssembly.GetLevel(tx) - 1
	}
	return &complexAssemblyImpl{
		Assembly:      newAssemblyImpl(tx, id, typ, buildDate, module, superAssembly),
		subAssemblies: gvstm.CreateTVar(newSmallSetImpl(tx)),
		level:         gvstm.CreateTVar(lvl),
	}
}

func (ca *complexAssemblyImpl) AddSubAssembly(tx Transaction, assembly Assembly) bool {
	return tx.Load(ca.subAssemblies).(smallSet).Add(tx, assembly)
}

func (ca *complexAssemblyImpl) RemoveSubAssembly(tx Transaction, assembly Assembly) bool {
	return tx.Load(ca.subAssemblies).(smallSet).Remove(tx, assembly)
}

func (ca *complexAssemblyImpl) GetSubAssemblies(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(tx.Load(ca.subAssemblies).(smallSet).ToSlice(tx))
}

func (ca *complexAssemblyImpl) GetLevel(tx Transaction) int {
	return tx.Load(ca.level).(int)
}

func (ca *complexAssemblyImpl) ClearPointers(tx Transaction) {
	ca.Assembly.ClearPointers(tx)
	tx.Store(ca.subAssemblies, nil)
	tx.Store(ca.level, -1)
}
