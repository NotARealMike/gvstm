package locks

import (
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
)

type compositePartImpl struct {
	DesignObj
	documentation Document
	usedIn        bag
	parts         LargeSet
	rootPart      AtomicPart
}

func newCompositePartImpl(tx Transaction, id int, typ string, buildDate int, documentation Document) CompositePart {
	cp := &compositePartImpl{
		DesignObj:     newDesignObjImpl(tx, id, typ, buildDate),
		documentation: documentation,
		usedIn:        newBagImpl(tx),
		parts:         beFactory.CreateLargeSet(tx),
		rootPart:      nil,
	}
	documentation.SetPart(tx, cp)
	return cp
}

func (cp *compositePartImpl) AddAssembly(tx Transaction, assembly BaseAssembly) {
	cp.usedIn.add(tx, assembly)
}

func (cp *compositePartImpl) AddPart(tx Transaction, part AtomicPart) bool {
	addedBefore := !cp.parts.Add(tx, part)
	if addedBefore {
		return false
	}
	part.SetCompositePart(tx, cp)
	if cp.rootPart == nil {
		cp.rootPart = part
	}
	return true
}

func (cp *compositePartImpl) SetRootPart(tx Transaction, part AtomicPart) {
	cp.rootPart = part
}

func (cp *compositePartImpl) GetRootPart(tx Transaction) AtomicPart {
	return cp.rootPart
}

func (cp *compositePartImpl) GetDocumentation(tx Transaction) Document {
	return cp.documentation
}

func (cp *compositePartImpl) GetParts(tx Transaction) LargeSet {
	return cp.parts
}

func (cp *compositePartImpl) RemoveAssembly(tx Transaction, assembly BaseAssembly) {
	cp.usedIn.remove(tx, assembly)
}

func (cp *compositePartImpl) GetUsedIn(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(cp.usedIn.toSlice(tx))
}
func (cp *compositePartImpl) ClearPointers(tx Transaction) {
	cp.documentation = nil
	cp.parts = nil
	cp.usedIn = nil
	cp.rootPart = nil
}
