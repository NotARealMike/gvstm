package gvstm

import (
    "gvstm/gvstm"
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type assemblyImpl struct {
    DesignObj
    superAssembly TVar
    module TVar
}

func newAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) Assembly {
    return &assemblyImpl{
        DesignObj:     newDesignObjImpl(tx, id, typ, buildDate),
        superAssembly: gvstm.CreateTVar(superAssembly),
        module:        gvstm.CreateTVar(module),
    }
}

func (a *assemblyImpl) GetSuperAssembly(tx Transaction) ComplexAssembly {
    return tx.Load(a.superAssembly).(ComplexAssembly)
}

func (a *assemblyImpl) GetModule(tx Transaction) Module {
    return tx.Load(a.module).(Module)
}

func (a *assemblyImpl) ClearPointers(tx Transaction) {
    tx.Store(a.superAssembly, nil)
    tx.Store(a.module, nil)
}
