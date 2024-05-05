package gvstm

import (
	"github.com/NotARealMike/gvstm/gvstm"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type designObjImpl struct {
	id        int
	typ       string
	buildDate TVar
}

func newDesignObjImpl(tx Transaction, id int, typ string, buildDate int) DesignObj {
	return &designObjImpl{
		id:        id,
		typ:       typ,
		buildDate: gvstm.CreateTVar(buildDate),
	}
}

func (do *designObjImpl) GetId(tx Transaction) int {
	return do.id
}

func (do *designObjImpl) GetBuildDate(tx Transaction) int {
	return tx.Load(do.buildDate).(int)
}

func (do *designObjImpl) GetType(tx Transaction) string {
	return do.typ
}

func (do *designObjImpl) UpdateBuildDate(tx Transaction) {
	bd := tx.Load(do.buildDate).(int)
	if bd%2 == 0 {
		tx.Store(do.buildDate, bd-1)
	} else {
		tx.Store(do.buildDate, bd+1)
	}
}

func (do *designObjImpl) NullOperation(tx Transaction) {}
