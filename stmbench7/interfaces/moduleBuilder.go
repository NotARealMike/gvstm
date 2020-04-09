package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
)

type ModuleBuilder struct {
    idPool IDPool
    manualBuilder *ManualBuilder
    assemblyBuilder *AssemblyBuilder
}

func NewModuleBuilder(tx Transaction, baseAssemblyIndex, complexAssemblyIndex Index) *ModuleBuilder {
    return &ModuleBuilder{
        idPool:          beFactory.CreateIDPool(tx, internal.NumModules),
        manualBuilder:   NewManualBuilder(tx),
        assemblyBuilder: NewAssemblyBuilder(tx, baseAssemblyIndex, complexAssemblyIndex),
    }
}

func (mb *ModuleBuilder) CreateAndRegisterModule(tx Transaction) (Module, *OpFailedError) {
    moduleID, err := mb.idPool.GetID(tx)
    if err != nil {
        return nil, err
    }

    manual, err := mb.manualBuilder.CreateManual(tx, moduleID)
    if err != nil {
        return nil, err
    }

    typ := createType()
    date := createBuildDate(internal.MinModuleDate, internal.MaxModuleDate)

    module := doFactory.CreateModule(tx, moduleID, typ, date, manual)
    designRoot, err := mb.assemblyBuilder.createAndRegisterComplexAssembly(tx, module, nil)
    if err != nil {
        return nil, err
    }
    module.SetDesignRoot(tx, designRoot)
    return module, nil
}

func (mb *ModuleBuilder) GetAssemblyBuilder() *AssemblyBuilder {
    return mb.assemblyBuilder
}
