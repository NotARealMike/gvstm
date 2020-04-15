package gvstm

import (
    "gvstm/gvstm"
    . "gvstm/stm"
)

type bagImpl struct {
    list TVar
}

func newBagImpl(tx Transaction) *bagImpl {
    return &bagImpl{
        list : gvstm.CreateTVar([]interface{}{}),
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
            index = 1
            break
        }
    }
    if index == -1 {
        return false
    }
    newList := make([]interface{}, len(oldList)-1)
    for i := range oldList {
        if i < index {
            newList[i] = oldList[i]
        } else if i > index {
            newList[i] = oldList[i+1]
        }
    }
    tx.Store(b.list, newList)
    return true
}

func (b *bagImpl) toSlice(tx Transaction) []interface{} {
    return tx.Load(b.list).([]interface{})
}
