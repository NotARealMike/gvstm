package operations

type OperationExecutor interface {
	Execute() (int, error)
}

type OperationExecutorFactory struct {
	CreateExecutor func(operation Operation) OperationExecutor
}

var OEFactory OperationExecutorFactory
