package locks

import (
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type atomicPartImpl struct {
	DesignObj
	x, y     int
	from, to smallSet
	partOf   CompositePart
}

func newAtomicPartImpl(tx Transaction, id int, typ string, buildDate int, x, y int) AtomicPart {
	return &atomicPartImpl{
		DesignObj: newDesignObjImpl(tx, id, typ, buildDate),
		x:         x,
		y:         y,
		from:      newSmallSetImpl(tx),
		to:        newSmallSetImpl(tx),
		partOf:    nil,
	}
}

func (ap *atomicPartImpl) ConnectTo(tx Transaction, destination AtomicPart, typ string, length int) {
	connection := NewConnectionImpl(ap, destination, typ, length)
	ap.to.Add(tx, connection)
	destination.AddConnectionFromOtherPart(tx, connection.GetReversed())
}

func (ap *atomicPartImpl) AddConnectionFromOtherPart(tx Transaction, connection Connection) {
	ap.from.Add(tx, connection)
}

func (ap *atomicPartImpl) SetCompositePart(tx Transaction, partOf CompositePart) {
	ap.partOf = partOf
}

func (ap *atomicPartImpl) GetNumToConnections(tx Transaction) int {
	return ap.to.Size(tx)
}

func (ap *atomicPartImpl) GetToConnections(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(ap.to.ToSlice(tx))
}

func (ap *atomicPartImpl) GetFromConnections(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(ap.from.ToSlice(tx))
}

func (ap *atomicPartImpl) GetPartOf(tx Transaction) CompositePart {
	return ap.partOf
}

func (ap *atomicPartImpl) SwapXY(tx Transaction) {
	tmp := ap.y
	ap.y = ap.x
	ap.x = tmp
}

func (ap *atomicPartImpl) GetX(tx Transaction) int {
	return ap.x
}

func (ap *atomicPartImpl) GetY(tx Transaction) int {
	return ap.y
}

func (ap *atomicPartImpl) ClearPointers(tx Transaction) {
	ap.x = 0
	ap.y = 0
	ap.to = nil
	ap.from = nil
	ap.partOf = nil
}
