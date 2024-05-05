package operations

import (
	. "github.com/NotARealMike/gvstm/stm"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"sync"
)

var (
	globalStructureWriteLock = &sync.RWMutex{}
	globalStructureReadLock  = globalStructureWriteLock.RLocker()

	assemblyWriteLocks []sync.Locker
	assemblyReadLocks  []sync.Locker

	compositePartWriteLock = &sync.RWMutex{}
	compositePartReadLock  = compositePartWriteLock.RLocker()
	atomicPartWriteLock    = &sync.RWMutex{}
	atomicPartReadLock     = atomicPartWriteLock.RLocker()
	documentWriteLock      = &sync.RWMutex{}
	documentReadLock       = documentWriteLock.RLocker()
	manualWriteLock        = &sync.RWMutex{}
	manualReadLock         = manualWriteLock.RLocker()
)

func init() {
	assemblyWriteLocks = make([]sync.Locker, internal.NumAssemblyLevels+1)
	assemblyReadLocks = make([]sync.Locker, internal.NumAssemblyLevels+1)
	for i := 1; i <= internal.NumAssemblyLevels; i++ {
		assemblyWriteLocks[i] = &sync.RWMutex{}
		assemblyReadLocks[i] = assemblyWriteLocks[i].(*sync.RWMutex).RLocker()
	}
}

type mgTransaction struct {
	isReadAcquired, isWriteAcquired []bool
}

func createMGTransaction(isReadAcquired, isWriteAcquired []bool) *mgTransaction {
	tx := &mgTransaction{
		isReadAcquired:  make([]bool, internal.NumAssemblyLevels+1),
		isWriteAcquired: make([]bool, internal.NumAssemblyLevels+1),
	}
	for i, v := range isReadAcquired {
		tx.isReadAcquired[i] = v
	}
	for i, v := range isWriteAcquired {
		tx.isWriteAcquired[i] = v
	}
	return tx
}

func (tx *mgTransaction) Load(tVar TVar) interface{} {
	panic("mgTransaction.Load should never be called!")
}

func (tx *mgTransaction) Store(tVar TVar, value interface{}) {
	panic("mgTransaction.Store should never be called!")
}

func ReadLockAssemblyLevel(tx Transaction, level int) {
	if tx == nil {
		return
	}
	if tx.(*mgTransaction).isReadAcquired[level] || tx.(*mgTransaction).isWriteAcquired[level] {
		return
	}
	assemblyReadLocks[level].Lock()
	tx.(*mgTransaction).isReadAcquired[level] = true
}

func WriteLockAssemblyLevel(tx Transaction, level int) {
	if tx == nil {
		return
	}
	if tx.(*mgTransaction).isWriteAcquired[level] {
		return
	}
	assemblyWriteLocks[level].Lock()
	tx.(*mgTransaction).isWriteAcquired[level] = true
}

type mgOperationExecutor struct {
	operation                                     Operation
	locksToAcquire                                []sync.Locker
	assReadLocksToAcquire, assWriteLocksToAcquire []bool
}

func CreateMGOperationExecutor(operation Operation) OperationExecutor {
	locksToAcquire := make([]sync.Locker, 0)
	assReadLocksToAcquire := make([]bool, internal.NumAssemblyLevels+1)
	assWriteLocksToAcquire := make([]bool, internal.NumAssemblyLevels+1)

	if operation.ID().OpType == StructureModification {
		return &mgOperationExecutor{
			operation:      operation,
			locksToAcquire: append(locksToAcquire, globalStructureWriteLock),
		}
	}

	locksToAcquire = append(locksToAcquire, globalStructureReadLock)

	switch operation.ID().Name {
	case "T1", "T6", "Q7", "ST1", "ST9", "OP1", "OP2", "OP3":
		locksToAcquire = append(locksToAcquire, atomicPartReadLock)
	case "T2a", "T2b", "T2c", "T3a", "T3b", "T3c", "ST6", "ST10", "OP9", "OP10", "OP15":
		locksToAcquire = append(locksToAcquire, atomicPartWriteLock)
	case "T4", "ST2":
		locksToAcquire = append(locksToAcquire, documentReadLock)
	case "T5", "ST7":
		locksToAcquire = append(locksToAcquire, documentWriteLock)
	case "Q6":
		locksToAcquire = append(locksToAcquire, assemblyReadLocks[1])
		assReadLocksToAcquire[1] = true
		locksToAcquire = append(locksToAcquire, compositePartReadLock)
		for i := internal.NumAssemblyLevels; i > 1; i-- {
			locksToAcquire = append(locksToAcquire, assemblyReadLocks[i])
			assReadLocksToAcquire[i] = true
		}
	case "ST3":
		locksToAcquire = append(locksToAcquire, assemblyReadLocks[1])
		assReadLocksToAcquire[1] = true
		for i := internal.NumAssemblyLevels; i > 1; i-- {
			locksToAcquire = append(locksToAcquire, assemblyReadLocks[i])
			assReadLocksToAcquire[i] = true
		}
	case "ST4":
		locksToAcquire = append(locksToAcquire, assemblyReadLocks[1])
		assReadLocksToAcquire[1] = true
		locksToAcquire = append(locksToAcquire, documentReadLock)
	case "ST5":
		locksToAcquire = append(locksToAcquire, assemblyReadLocks[1])
		assReadLocksToAcquire[1] = true
		locksToAcquire = append(locksToAcquire, compositePartReadLock)
	case "ST8":
		locksToAcquire = append(locksToAcquire, assemblyWriteLocks[1])
		assWriteLocksToAcquire[1] = true
		for i := internal.NumAssemblyLevels; i > 1; i-- {
			locksToAcquire = append(locksToAcquire, assemblyWriteLocks[i])
			assWriteLocksToAcquire[i] = true
		}
	case "OP4", "OP5":
		locksToAcquire = append(locksToAcquire, manualReadLock)
	case "OP6", "OP12":
		// none
	case "OP7":
		locksToAcquire = append(locksToAcquire, assemblyReadLocks[1])
		assReadLocksToAcquire[1] = true
	case "OP8":
		locksToAcquire = append(locksToAcquire, compositePartReadLock)
	case "OP11":
		locksToAcquire = append(locksToAcquire, manualWriteLock)
	case "OP13":
		locksToAcquire = append(locksToAcquire, assemblyWriteLocks[1])
		assWriteLocksToAcquire[1] = true
	case "OP14":
		locksToAcquire = append(locksToAcquire, compositePartWriteLock)
	default:
		panic("Unknown operation type: " + operation.ID().Name)
	}

	return &mgOperationExecutor{
		operation:              operation,
		locksToAcquire:         locksToAcquire,
		assReadLocksToAcquire:  assReadLocksToAcquire,
		assWriteLocksToAcquire: assWriteLocksToAcquire,
	}
}

func (e *mgOperationExecutor) Execute() (result int, err error) {
	tx := createMGTransaction(e.assReadLocksToAcquire, e.assWriteLocksToAcquire)
	for _, lock := range e.locksToAcquire {
		lock.Lock()
	}
	defer func() {
		for i := 1; i <= internal.NumAssemblyLevels; i++ {
			if tx.isReadAcquired[i] && !e.assReadLocksToAcquire[i] {
				assemblyReadLocks[i].Unlock()
			}
			if tx.isWriteAcquired[i] && !e.assWriteLocksToAcquire[i] {
				assemblyWriteLocks[i].Unlock()
			}
		}

		for _, lock := range e.locksToAcquire {
			lock.Unlock()
		}
	}()

	return e.operation.Perform(tx)
}
