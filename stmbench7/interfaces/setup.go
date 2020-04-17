package interfaces

import (
    "fmt"
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
    "math/rand"
)

type Setup interface {
    Module() Module

    AtomicPartIDIndex() Index
    AtomicPartBuildDateIndex() Index
    CompositePartIDIndex() Index
    BaseAssemblyIDIndex() Index
    ComplexAssemblyIDIndex() Index
    DocumentTitleIndex() Index

    CompositePartBuilder() CompositePartBuilder
    AssemblyBuilder() AssemblyBuilder
    modBuilder() moduleBuilder
}

type setupImpl struct {
    module Module

    atomicPartIDIndex,
    atomicPartBuildDateIndex,
    compositePartIDIndex,
    baseAssemblyIDIndex,
    complexAssemblyIDIndex,
    documentTitleIndex Index

    compositePartBuilder CompositePartBuilder
    moduleBuilder moduleBuilder
}

func (s *setupImpl) Module() Module {
    return s.module
}

func (s *setupImpl) AtomicPartIDIndex() Index {
    return s.atomicPartIDIndex
}

func (s *setupImpl) AtomicPartBuildDateIndex() Index {
    return s.atomicPartBuildDateIndex
}

func (s *setupImpl) CompositePartIDIndex() Index {
    return s.compositePartIDIndex
}

func (s *setupImpl) BaseAssemblyIDIndex() Index {
    return s.baseAssemblyIDIndex
}

func (s *setupImpl) ComplexAssemblyIDIndex() Index {
    return s.complexAssemblyIDIndex
}

func (s *setupImpl) DocumentTitleIndex() Index {
    return s.documentTitleIndex
}

func (s *setupImpl) CompositePartBuilder() CompositePartBuilder {
    return s.compositePartBuilder
}

func (s *setupImpl) AssemblyBuilder() AssemblyBuilder {
    return s.moduleBuilder.getAssemblyBuilder()
}

func (s *setupImpl) modBuilder() moduleBuilder {
    return s.moduleBuilder
}

func NewSetup(tx Transaction) Setup {
    s := &setupImpl{
        module:                   nil,
        atomicPartIDIndex:        BEFactory.CreateIndex(tx),
        atomicPartBuildDateIndex: BEFactory.CreateIndex(tx),
        compositePartIDIndex:     BEFactory.CreateIndex(tx),
        baseAssemblyIDIndex:      BEFactory.CreateIndex(tx),
        complexAssemblyIDIndex:   BEFactory.CreateIndex(tx),
        documentTitleIndex:       BEFactory.CreateIndex(tx),
        compositePartBuilder:     nil,
        moduleBuilder:            nil,
    }
    s.compositePartBuilder = newCompositePartBuilder(tx, s.compositePartIDIndex, s.documentTitleIndex, s.atomicPartIDIndex, s.atomicPartBuildDateIndex)
    s.moduleBuilder = newModuleBuilder(tx, s.baseAssemblyIDIndex, s.complexAssemblyIDIndex)
    s.module = setUpDataStructure(tx, s)
    return s
}

func setUpDataStructure(tx Transaction, setup Setup) Module {
    fmt.Println("Setting up the design library...")
    designLibrary := make([]CompositePart, internal.InitialTotalCompParts)
    compPartBuilder := setup.CompositePartBuilder()
    for i := 0; i < internal.InitialTotalCompParts; i++ {
        fmt.Printf("Component %d of %d\n", i+1, internal.InitialTotalCompParts)
        compPart, err := compPartBuilder.CreateAndRegisterCompositePart(tx)
        if err != nil {
            panic("Unexpected failure when creating composite part: " + err.Error())
        }
        designLibrary[i] = compPart
    }
    fmt.Println()

    fmt.Println("Setting up the module...")
    module, err := setup.modBuilder().createAndRegisterModule(tx)
    if err != nil {
        panic("Unexpected failure when creating module: " + err.Error())
    }
    for i, baseAssembly := range setup.BaseAssemblyIDIndex().GetRange(tx, internal.MinAssDate, internal.MaxAssDate) {
        fmt.Printf("Base assembly %d of %d\n", i+1, internal.InitialTotalBaseAssemblies)
        for i := 0; i < internal.NumCompPartPerAss; i++ {
            // TODO: Java version uses ThreadRandom.
            compositePartNumber := rand.Int()
            baseAssembly.(BaseAssembly).AddComponent(tx, designLibrary[compositePartNumber])
        }
    }
    fmt.Println()

    return module
}
