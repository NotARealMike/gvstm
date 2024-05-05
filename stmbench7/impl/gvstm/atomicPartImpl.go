package gvstm

import (
	"github.com/NotARealMike/gvstm/gvstm"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type atomicPartImpl struct {
	DesignObj
	x, y     TVar
	from, to TVar
	partOf   TVar
}

func newAtomicPartImpl(tx Transaction, id int, typ string, buildDate int, x, y int) AtomicPart {
	return &atomicPartImpl{
		DesignObj: newDesignObjImpl(tx, id, typ, buildDate),
		x:         gvstm.CreateTVar(x),
		y:         gvstm.CreateTVar(y),
		from:      gvstm.CreateTVar(newSmallSetImpl(tx)),
		to:        gvstm.CreateTVar(newSmallSetImpl(tx)),
		partOf:    gvstm.CreateTVar(nil),
	}
}

func (ap *atomicPartImpl) ConnectTo(tx Transaction, destination AtomicPart, typ string, length int) {
	connection := NewConnectionImpl(ap, destination, typ, length)
	tx.Load(ap.to).(smallSet).Add(tx, connection)
	destination.AddConnectionFromOtherPart(tx, connection.GetReversed())
}

func (ap *atomicPartImpl) AddConnectionFromOtherPart(tx Transaction, connection Connection) {
	tx.Load(ap.from).(smallSet).Add(tx, connection)
}

func (ap *atomicPartImpl) SetCompositePart(tx Transaction, partOf CompositePart) {
	tx.Store(ap.partOf, partOf)
}

func (ap *atomicPartImpl) GetNumToConnections(tx Transaction) int {
	return tx.Load(ap.to).(smallSet).Size(tx)
}

func (ap *atomicPartImpl) GetToConnections(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(tx.Load(ap.to).(smallSet).ToSlice(tx))
}

func (ap *atomicPartImpl) GetFromConnections(tx Transaction) ImmutableCollection {
	return NewImmutableCollectionImpl(tx.Load(ap.from).(smallSet).ToSlice(tx))
}

func (ap *atomicPartImpl) GetPartOf(tx Transaction) CompositePart {
	return tx.Load(ap.partOf).(CompositePart)
}

func (ap *atomicPartImpl) SwapXY(tx Transaction) {
	tmp := tx.Load(ap.y)
	tx.Store(ap.y, tx.Load(ap.x))
	tx.Store(ap.x, tmp)
}

func (ap *atomicPartImpl) GetX(tx Transaction) int {
	return tx.Load(ap.x).(int)
}

func (ap *atomicPartImpl) GetY(tx Transaction) int {
	return tx.Load(ap.y).(int)
}

func (ap *atomicPartImpl) ClearPointers(tx Transaction) {
	tx.Store(ap.x, 0)
	tx.Store(ap.y, 0)
	tx.Store(ap.to, nil)
	tx.Store(ap.from, nil)
	tx.Store(ap.partOf, nil)
}
