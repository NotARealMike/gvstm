package gvstm

import (
	"github.com/NotARealMike/gvstm/gvstm"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type compositePartImpl struct {
	DesignObj
	documentation TVar
	usedIn        TVar
	parts         TVar
	rootPart      TVar
}

func newCompositePartImpl(tx Transaction, id int, typ string, buildDate int, documentation Document) CompositePart {
	cp := &compositePartImpl{
		DesignObj:     newDesignObjImpl(tx, id, typ, buildDate),
		documentation: gvstm.CreateTVar(documentation),
		usedIn:        gvstm.CreateTVar(newBagImpl(tx)),
		parts:         gvstm.CreateTVar(beFactory.CreateLargeSet(tx)),
		rootPart:      gvstm.CreateTVar(nil),
	}
	documentation.SetPart(tx, cp)
	return cp
}

func (cp *compositePartImpl) AddAssembly(tx Transaction, assembly BaseAssembly) {
	tx.Load(cp.usedIn).(bag).add(tx, assembly)
}

func (cp *compositePartImpl) AddPart(tx Transaction, part AtomicPart) bool {
	addedBefore := !tx.Load(cp.parts).(LargeSet).Add(tx, part)
	if addedBefore {
		return false
	}
	part.SetCompositePart(tx, cp)
	if tx.Load(cp.rootPart) == nil {
		tx.Store(cp.rootPart, part)
	}
	return true
}

func (cp *compositePartImpl) SetRootPart(tx Transaction, part AtomicPart) {
	tx.Store(cp.rootPart, part)
}

func (cp *compositePartImpl) GetRootPart(tx Transaction) AtomicPart {
	return tx.Load(cp.rootPart).(AtomicPart)
}

func (cp *compositePartImpl) GetDocumentation(tx Transaction) Document {
	return tx.Load(cp.documentation).(Document)
}

func (cp *compositePartImpl) GetParts(tx Transaction) LargeSet {
	return tx.Load(cp.parts).(LargeSet)
}

func (cp *compositePartImpl) RemoveAssembly(tx Transaction, assembly BaseAssembly) {
	tx.Load(cp.usedIn).(bag).remove(tx, assembly)
}

func (cp *compositePartImpl) GetUsedIn(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(tx.Load(cp.usedIn).(bag).toSlice(tx))
}
func (cp *compositePartImpl) ClearPointers(tx Transaction) {
	tx.Store(cp.documentation, nil)
	tx.Store(cp.parts, nil)
	tx.Store(cp.usedIn, nil)
	tx.Store(cp.rootPart, nil)
}
