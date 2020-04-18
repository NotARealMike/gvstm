package locks

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type assemblyImpl struct {
    DesignObj
    superAssembly ComplexAssembly
    module Module
}

func newAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) Assembly {
    return &assemblyImpl{
        DesignObj:     newDesignObjImpl(tx, id, typ, buildDate),
        superAssembly: superAssembly,
        module:        module,
    }
}

func (a *assemblyImpl) GetSuperAssembly(tx Transaction) ComplexAssembly {
    return a.superAssembly
}

func (a *assemblyImpl) GetModule(tx Transaction) Module {
    return a.module
}

func (a *assemblyImpl) ClearPointers(tx Transaction) {
    a.superAssembly = nil
    a.module = nil
}
