package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
)

type moduleBuilder interface {
    createAndRegisterModule(tx Transaction) (Module, OpFailedError)
    getAssemblyBuilder() AssemblyBuilder
}

type moduleBuilderImpl struct {
    idPool IDPool
    manualBuilder manualBuilder
    assemblyBuilder AssemblyBuilder
}

func newModuleBuilder(tx Transaction, baseAssemblyIndex, complexAssemblyIndex Index) moduleBuilder {
    return &moduleBuilderImpl{
        idPool:          BEFactory.CreateIDPool(tx, internal.NumModules),
        manualBuilder:   newManualBuilder(tx),
        assemblyBuilder: newAssemblyBuilder(tx, baseAssemblyIndex, complexAssemblyIndex),
    }
}

func (mb *moduleBuilderImpl) createAndRegisterModule(tx Transaction) (Module, OpFailedError) {
    moduleID, err := mb.idPool.GetID(tx)
    if err != nil {
        return nil, NewOpFailedError("creating module: " + err.Error())
    }

    manual, err := mb.manualBuilder.createManual(tx, moduleID)
    if err != nil {
        return nil, NewOpFailedError("creating manual: " + err.Error())
    }

    typ := createType()
    date := createBuildDate(internal.MinModuleDate, internal.MaxModuleDate)

    module := DOFactory.CreateModule(tx, moduleID, typ, date, manual)
    designRoot, err := mb.assemblyBuilder.CreateAndRegisterAssembly(tx, module, nil)
    if err != nil {
        return nil, NewOpFailedError("creating design root: " + err.Error())
    }
    module.SetDesignRoot(tx, designRoot.(ComplexAssembly))
    return module, nil
}

func (mb *moduleBuilderImpl) getAssemblyBuilder() AssemblyBuilder {
    return mb.assemblyBuilder
}
