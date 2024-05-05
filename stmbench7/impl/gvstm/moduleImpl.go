package gvstm

import (
	"github.com/NotARealMike/gvstm/gvstm"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
)

type moduleImpl struct {
	DesignObj
	man        Manual
	designRoot TVar
}

func newModuleImpl(tx Transaction, id int, typ string, buildDate int, man Manual) Module {
	m := &moduleImpl{
		DesignObj:  newDesignObjImpl(tx, id, typ, buildDate),
		man:        man,
		designRoot: gvstm.CreateTVar(nil),
	}
	man.SetModule(tx, m)
	return m
}

func (m *moduleImpl) SetDesignRoot(tx Transaction, designRoot ComplexAssembly) {
	tx.Store(m.designRoot, designRoot)
}

func (m *moduleImpl) GetDesignRoot(tx Transaction) ComplexAssembly {
	return tx.Load(m.designRoot).(ComplexAssembly)
}

func (m *moduleImpl) GetManual(tx Transaction) Manual {
	return m.man
}
