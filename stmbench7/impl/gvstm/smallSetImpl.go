package gvstm

import (
    "gvstm/gvstm"
    . "gvstm/stm"
)

type smallSet interface {
    add(tx Transaction, element interface{}) bool
    remove(tx Transaction, element interface{}) bool
    contains(tx Transaction, element interface{}) bool
    size(tx Transaction) int
    toSlice(tx Transaction) []interface{}
}

type smallSetImpl struct {
    set TVar
}

func newSmallSetImpl(tx Transaction) smallSet {
    return &smallSetImpl{
        set : gvstm.CreateTVar([]interface{}{}),
    }
}

func (s *smallSetImpl) add(tx Transaction, element interface{}) bool {
    if s.contains(tx, element) {
        return false
    }
    oldSet := tx.Load(s.set).([]interface{})
    newSet := make([]interface{}, len(oldSet)+1)
    for i := range oldSet {
        newSet[i] = oldSet[i]
    }
    newSet[len(oldSet)] = element
    tx.Store(s.set, newSet)
    return true
}

func (s *smallSetImpl) remove(tx Transaction, element interface{}) bool {
    if !s.contains(tx, element) {
        return false
    }
    oldSet := tx.Load(s.set).([]interface{})
    newSet := make([]interface{}, len(oldSet)-1)
    found := false
    for i := range oldSet {
        if oldSet[i] == element {
            found = true
        }
        if !found {
            newSet[i] = oldSet[i]
        } else {
            newSet[i] = oldSet[i+1]
        }
    }
    tx.Store(s.set, newSet)
    return true
}

func (s *smallSetImpl) contains(tx Transaction, element interface{}) bool {
    set := tx.Load(s.set).([]interface{})
    for _, e := range set {
        if e == element {
            return true
        }
    }
    return false
}

func (s *smallSetImpl) size(tx Transaction) int {
    return len(tx.Load(s.set).([]interface{}))
}

func (s *smallSetImpl) toSlice(tx Transaction) []interface{} {
    return tx.Load(s.set).([]interface{})
}