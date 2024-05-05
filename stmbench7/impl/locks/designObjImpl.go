package locks

import (
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type designObjImpl struct {
	id        int
	typ       string
	buildDate int
}

func newDesignObjImpl(tx Transaction, id int, typ string, buildDate int) DesignObj {
	return &designObjImpl{
		id:        id,
		typ:       typ,
		buildDate: buildDate,
	}
}

func (do *designObjImpl) GetId(tx Transaction) int {
	return do.id
}

func (do *designObjImpl) GetBuildDate(tx Transaction) int {
	return do.buildDate
}

func (do *designObjImpl) GetType(tx Transaction) string {
	return do.typ
}

func (do *designObjImpl) UpdateBuildDate(tx Transaction) {
	if do.buildDate%2 == 0 {
		do.buildDate--
	} else {
		do.buildDate++
	}
}

func (do *designObjImpl) NullOperation(tx Transaction) {}
