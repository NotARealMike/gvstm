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
    offset := 0
    for i := range newSet {
        if s.set[i] == element {
            offset = 1
        }
        newSet[i] = s.set[i+offset]
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