package locks

import (
    . "gvstm/stm"
    "gvstm/stmbench7/interfaces"
)

type smallSet interface {
    interfaces.LargeSet
}

type smallSetImpl struct {
    interfaces.LargeSet
}

func newSmallSetImpl(tx Transaction) smallSet {
    return &smallSetImpl{
        newLargeSetImpl(tx),
    }
}
