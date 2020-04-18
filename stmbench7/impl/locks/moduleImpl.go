package locks

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type moduleImpl struct {
    DesignObj
    man Manual
    designRoot ComplexAssembly
}

func newModuleImpl(tx Transaction, id int, typ string, buildDate int, man Manual) Module {
    m := &moduleImpl{
        DesignObj:  newDesignObjImpl(tx, id, typ, buildDate),
        man:        man,
        designRoot: nil,
    }
    man.SetModule(tx, m)
    return m
}

func (m *moduleImpl) SetDesignRoot(tx Transaction, designRoot ComplexAssembly) {
    m.designRoot = designRoot
}

func (m *moduleImpl) GetDesignRoot(tx Transaction) ComplexAssembly {
    return m.designRoot
}

func (m *moduleImpl) GetManual(tx Transaction) Manual {
    return m.man
}