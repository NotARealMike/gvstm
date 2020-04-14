package gvstm

import . "gvstm/stm"

// TODO: stub
type smallSetImpl struct {}

func newSmallSetImpl(tx Transaction) *smallSetImpl {
    return &smallSetImpl{}
}

func (s *smallSetImpl) add(tx Transaction, element interface{}) bool {
    return false
}

func (s *smallSetImpl) remove(tx Transaction, element interface{}) bool {
    return false
}

func (s *smallSetImpl) contains(tx Transaction, element interface{}) bool {
    return false
}

func (s *smallSetImpl) size(tx Transaction) int {
    return 0
}

func (s *smallSetImpl) toSlice(tx Transaction) []interface{} {
    return nil
}