package interfaces

import . "gvstm/stm"

// A factory for creating the benchmark data structures.
// STM implementations should populate an instance with their respective factory methods.
type DesignObjFactory struct {
    CreateAtomicPart func(tx Transaction, id int, typ string, buildDate int, x, y int) AtomicPart
    CreateConnection func(from, to AtomicPart, typ string, length int) Connection
    CreateBaseAssembly func(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) BaseAssembly
    CreateComplexAssembly func(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) ComplexAssembly
    CreateCompositePart func(tx Transaction, id int, typ string, buildDate int, documentation Document) CompositePart
    CreateDocument func(tx Transaction, id int, title, text string) Document
    CreateManual func(tx Transaction, id int, title, text string) Manual
    CreateModule func(tx Transaction, id int, typ string, buildDate int, man Manual) Module
}

// Holds the factory methods used to create the backend data structures.
type BackendFactory struct {
    CreateLargeSet func(tx Transaction) LargeSet
    CreateIndex func(tx Transaction) Index
    CreateIDPool func(tx Transaction, maxNumberOfIDs int) IDPool
}

var (
    doFactory DesignObjFactory
    beFactory BackendFactory
)

func SetFactories(designObjectFactory DesignObjFactory, backendFactory BackendFactory) {
    doFactory = designObjectFactory
    beFactory = backendFactory
}