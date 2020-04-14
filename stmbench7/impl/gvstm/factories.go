package gvstm

import "gvstm/stmbench7/interfaces"

var (
    doFactory = interfaces.DesignObjFactory{
        CreateAtomicPart:      newAtomicPartImpl,
        CreateConnection:      interfaces.NewConnectionImpl,
        CreateBaseAssembly:    newBaseAssemblyImpl,
        CreateComplexAssembly: newComplexAssemblyImpl,
        CreateCompositePart:   newCompositePartImpl,
        CreateDocument:        newDocumentImpl,
        CreateManual:          newManualImpl,
        CreateModule:          newModuleImpl,
    }
    beFactory = interfaces.BackendFactory{
        CreateLargeSet: nil,
        CreateIndex:    nil,
        CreateIDPool:   nil,
    }
)
