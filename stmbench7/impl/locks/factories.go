package locks

import "gvstm/stmbench7/interfaces"

var (
	CGLocksInitialiser = interfaces.SynchMethodInitialiser{
		DOFactory: doFactory,
		BEFactory: beFactory,
	}
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
		CreateLargeSet: newLargeSetImpl,
		CreateIndex:    newIndexImpl,
		CreateIDPool:   newIDPoolImpl,
	}
	MGLocksInitialiser = interfaces.SynchMethodInitialiser{
		DOFactory: mgDOFactory,
		BEFactory: beFactory,
	}
	mgDOFactory = interfaces.DesignObjFactory{
		CreateAtomicPart:      newAtomicPartImpl,
		CreateConnection:      interfaces.NewConnectionImpl,
		CreateBaseAssembly:    newBaseAssemblyImpl,
		CreateComplexAssembly: newMGComplexAssemblyImpl,
		CreateCompositePart:   newCompositePartImpl,
		CreateDocument:        newDocumentImpl,
		CreateManual:          newManualImpl,
		CreateModule:          newModuleImpl,
	}
)
