package gvstm

import (
	. "github.com/NotARealMike/gvstm/stm"
	"github.com/NotARealMike/gvstm/stmbench7/interfaces"
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
