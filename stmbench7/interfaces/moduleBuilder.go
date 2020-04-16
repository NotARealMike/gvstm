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
        idPool:          beFactory.CreateIDPool(tx, internal.NumModules),
        manualBuilder:   newManualBuilder(tx),
        assemblyBuilder: newAssemblyBuilder(tx, baseAssemblyIndex, complexAssemblyIndex),
    }
}

func (mb *moduleBuilderImpl) createAndRegisterModule(tx Transaction) (Module, OpFailedError) {
    moduleID, err := mb.idPool.GetID(tx)
    if err != nil {
        return nil, err
    }

    manual, err := mb.manualBuilder.createManual(tx, moduleID)
    if err != nil {
        return nil, err
    }

    typ := createType()
    date := createBuildDate(internal.MinModuleDate, internal.MaxModuleDate)

    module := doFactory.CreateModule(tx, moduleID, typ, date, manual)
    designRoot, err := mb.assemblyBuilder.CreateAndRegisterAssembly(tx, module, nil)
    if err != nil {
        return nil, err
    }
    module.SetDesignRoot(tx, designRoot.(ComplexAssembly))
    return module, nil
}

func (mb *moduleBuilderImpl) getAssemblyBuilder() AssemblyBuilder {
    return mb.assemblyBuilder
}
