package operations

import "sync"

var (
    globalWriteLock = &sync.RWMutex{}
    globalReadLock = globalWriteLock.RLocker()
)

type cgOperationExecutor struct {
    operation Operation
    lock sync.Locker
}

func CreateCGOperationExecutor(operation Operation) OperationExecutor {
    var lock sync.Locker
    switch operation.ID().OpType {
    case OperationRO, TraversalRO, ShortTraversalRO:
        lock = globalReadLock
    default:
        lock = globalWriteLock
    }
    return &cgOperationExecutor{operation, lock}
}

func (e *cgOperationExecutor) Execute() (result int, err error) {
    e.lock.Lock()
    defer e.lock.Unlock()
    return e.operation.Perform(nil)
}
