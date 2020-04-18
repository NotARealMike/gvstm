package locks

import (
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
    set []interface{}
}

func newSmallSetImpl(tx Transaction) smallSet {
    return &smallSetImpl{
        set : []interface{}{},
    }
}

func (s *smallSetImpl) add(tx Transaction, element interface{}) bool {
    if s.contains(tx, element) {
        return false
    }
    newSet := make([]interface{}, len(s.set)+1)
    for i := range s.set {
        newSet[i] = s.set[i]
    }
    newSet[len(s.set)] = element
    s.set = newSet
    return true
}

func (s *smallSetImpl) remove(tx Transaction, element interface{}) bool {
    if !s.contains(tx, element) {
        return false
    }
    newSet := make([]interface{}, len(s.set)-1)
    found := false
    for i := range s.set {
        if s.set[i] == element {
            found = true
        }
        if !found {
            newSet[i] = s.set[i]
        } else {
            newSet[i] = s.set[i+1]
        }
    }
    s.set = newSet
    return true
}

func (s *smallSetImpl) contains(tx Transaction, element interface{}) bool {
    for _, e := range s.set {
        if e == element {
            return true
        }
    }
    return false
}

func (s *smallSetImpl) size(tx Transaction) int {
    return len(s.set)
}

func (s *smallSetImpl) toSlice(tx Transaction) []interface{} {
    return s.set
}