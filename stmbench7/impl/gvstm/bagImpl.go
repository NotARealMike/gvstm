package gvstm

import . "gvstm/stm"

// TODO: stub
type bagImpl struct {}

func newBagImpl(tx Transaction) *bagImpl {
    return &bagImpl{}
}

func (b *bagImpl) add(tx Transaction, element interface{}) {}

func (b *bagImpl) remove(tx Transaction, element interface{}) bool {
    return false
}

func (b *bagImpl) toSlice(tx Transaction) []interface{} {
    return nil
}