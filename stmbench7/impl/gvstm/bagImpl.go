package gvstm

import (
	"gvstm/gvstm"
	. "gvstm/stm"
)

type bag interface {
	add(tx Transaction, element interface{})
	remove(tx Transaction, element interface{}) bool
	toSlice(tx Transaction) []interface{}
}

type bagImpl struct {
	list TVar
}

func newBagImpl(tx Transaction) bag {
	return &bagImpl{
		list: gvstm.CreateTVar([]interface{}{}),
	}
}

func (b *bagImpl) add(tx Transaction, element interface{}) {
	oldList := tx.Load(b.list).([]interface{})
	newList := make([]interface{}, len(oldList)+1)
	for i := range oldList {
		newList[i] = oldList[i]
	}
	newList[len(oldList)] = element
	tx.Store(b.list, newList)
}

func (b *bagImpl) remove(tx Transaction, element interface{}) bool {
	oldList := tx.Load(b.list).([]interface{})
	index := -1
	for i := range oldList {
		if oldList[i] == element {
			index = i
			break
		}
	}
	if index == -1 {
		return false
	}
	newList := make([]interface{}, len(oldList)-1)
	for i := 0; i < index; i++ {
		newList[i] = oldList[i]
	}
	for i := index; i < len(newList); i++ {
		newList[i] = oldList[i+1]
	}
	tx.Store(b.list, newList)
	return true
}

func (b *bagImpl) toSlice(tx Transaction) []interface{} {
	return tx.Load(b.list).([]interface{})
}
