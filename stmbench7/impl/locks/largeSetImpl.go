package locks

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type largeSetImpl struct {
    set []interface{}
}

func newLargeSetImpl(tx Transaction) LargeSet {
    return &largeSetImpl{
        set : []interface{}{},
    }
}

func (s *largeSetImpl) Add(tx Transaction, element interface{}) bool {
    if s.Contains(tx, element) {
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

func (s *largeSetImpl) Remove(tx Transaction, element interface{}) bool {
    if !s.Contains(tx, element) {
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

func (s *largeSetImpl) Contains(tx Transaction, element interface{}) bool {
    for _, e := range s.set {
        if e == element {
            return true
        }
    }
    return false
}

func (s *largeSetImpl) Size(tx Transaction) int {
    return len(s.set)
}

func (s *largeSetImpl) ToSlice(tx Transaction) []interface{} {
    return s.set
}