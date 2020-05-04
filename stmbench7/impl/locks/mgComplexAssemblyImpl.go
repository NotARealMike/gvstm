package locks

import (
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
	"gvstm/stmbench7/operations"
)

type mgComplexAssemblyImpl struct {
	ComplexAssembly
}

func newMGComplexAssemblyImpl(tx Transaction, id int, typ string, buildDate int, module Module, superAssembly ComplexAssembly) ComplexAssembly {
	return &mgComplexAssemblyImpl{
		ComplexAssembly: newComplexAssemblyImpl(tx, id, typ, buildDate, module, superAssembly),
	}
}

func (ca *mgComplexAssemblyImpl) GetBuildDate(tx Transaction) int {
	operations.ReadLockAssemblyLevel(tx, ca.GetLevel(tx))
	return ca.ComplexAssembly.GetBuildDate(tx)
}

func (ca *mgComplexAssemblyImpl) GetType(tx Transaction) string {
	operations.ReadLockAssemblyLevel(tx, ca.GetLevel(tx))
	return ca.ComplexAssembly.GetType(tx)
}

func (ca *mgComplexAssemblyImpl) UpdateBuildDate(tx Transaction) {
	operations.WriteLockAssemblyLevel(tx, ca.GetLevel(tx))
	ca.ComplexAssembly.UpdateBuildDate(tx)
}

func (ca *mgComplexAssemblyImpl) NullOperation(tx Transaction) {
	operations.ReadLockAssemblyLevel(tx, ca.GetLevel(tx))
}
