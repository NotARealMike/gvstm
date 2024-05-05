package gvstm

import (
	"github.com/NotARealMike/gvstm/gvstm"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type largeSetImpl struct {
	set TVar
}

func newLargeSetImpl(tx Transaction) LargeSet {
	return &largeSetImpl{
		set: gvstm.CreateTVar([]interface{}{}),
	}
}

func (s *largeSetImpl) Add(tx Transaction, element interface{}) bool {
	if s.Contains(tx, element) {
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

func (s *largeSetImpl) Remove(tx Transaction, element interface{}) bool {
	if !s.Contains(tx, element) {
		return false
	}
	oldSet := tx.Load(s.set).([]interface{})
	newSet := make([]interface{}, len(oldSet)-1)
	offset := 0
	for i := range newSet {
		if oldSet[i] == element {
			offset = 1
		}
		newSet[i] = oldSet[i+offset]
	}
	tx.Store(s.set, newSet)
	return true
}

func (s *largeSetImpl) Contains(tx Transaction, element interface{}) bool {
	set := tx.Load(s.set).([]interface{})
	for _, e := range set {
		if e == element {
			return true
		}
	}
	return false
}

func (s *largeSetImpl) Size(tx Transaction) int {
	return len(tx.Load(s.set).([]interface{}))
}

func (s *largeSetImpl) ToSlice(tx Transaction) []interface{} {
	return tx.Load(s.set).([]interface{})
}
