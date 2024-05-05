package operations

import (
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"math/rand"
)

func createOperation1(setup Setup) Operation {
	atomicPartOperation := func(tx Transaction, part AtomicPart) {
		part.NullOperation(tx)
	}
	return &operationImpl{
		id: OP1,
		f:  createOps1_9_15(setup, atomicPartOperation),
	}
}

func createOperation2(setup Setup) Operation {
	atomicPartOperation := func(tx Transaction, part AtomicPart) {
		part.NullOperation(tx)
	}
	return &operationImpl{
		id: OP2,
		f:  createOps2_3_10(setup, 1, atomicPartOperation),
	}
}

func createOperation3(setup Setup) Operation {
	atomicPartOperation := func(tx Transaction, part AtomicPart) {
		part.NullOperation(tx)
	}
	return &operationImpl{
		id: OP3,
		f:  createOps2_3_10(setup, 10, atomicPartOperation),
	}
}

func createOperation4(setup Setup) Operation {
	f := func(tx Transaction) (int, error) {
		manual := setup.Module().GetManual(tx)
		return manual.CountOccurences(tx, 'I'), nil
	}
	return &operationImpl{
		id: OP4,
		f:  f,
	}
}

func createOperation5(setup Setup) Operation {
	f := func(tx Transaction) (int, error) {
		manual := setup.Module().GetManual(tx)
		return manual.CheckFirstLastCharTheSame(tx), nil
	}
	return &operationImpl{
		id: OP5,
		f:  f,
	}
}

func createOperation6(setup Setup) Operation {
	complexAssemblyOperation := func(tx Transaction, assembly ComplexAssembly) {
		assembly.NullOperation(tx)
	}
	return &operationImpl{
		id: OP6,
		f:  createOps6_12(setup, complexAssemblyOperation),
	}
}

func createOperation7(setup Setup) Operation {
	baseAssemblyOperation := func(tx Transaction, assembly BaseAssembly) {
		assembly.NullOperation(tx)
	}
	return &operationImpl{
		id: OP7,
		f:  createOps7_13(setup, baseAssemblyOperation),
	}
}

func createOperation8(setup Setup) Operation {
	componentOperation := func(tx Transaction, component CompositePart) {
		component.NullOperation(tx)
	}
	return &operationImpl{
		id: OP8,
		f:  createOps8_14(setup, componentOperation),
	}
}

func createOperation9(setup Setup) Operation {
	atomicPartOperation := func(tx Transaction, part AtomicPart) {
		part.SwapXY(tx)
	}
	return &operationImpl{
		id: OP9,
		f:  createOps1_9_15(setup, atomicPartOperation),
	}
}

func createOperation10(setup Setup) Operation {
	atomicPartOperation := func(tx Transaction, part AtomicPart) {
		part.SwapXY(tx)
	}
	return &operationImpl{
		id: OP10,
		f:  createOps2_3_10(setup, 1, atomicPartOperation),
	}
}

func createOperation11(setup Setup) Operation {
	f := func(tx Transaction) (int, error) {
		manual := setup.Module().GetManual(tx)
		if manual.StartsWith(tx, 'i') {
			return manual.ReplaceChar(tx, 'i', 'I'), nil
		} else if manual.StartsWith(tx, 'I') {
			return manual.ReplaceChar(tx, 'I', 'i'), nil
		}
		panic("OP11: unexpected Manual.text!")
	}
	return &operationImpl{
		id: OP11,
		f:  f,
	}
}

func createOperation12(setup Setup) Operation {
	complexAssemblyOperation := func(tx Transaction, assembly ComplexAssembly) {
		assembly.UpdateBuildDate(tx)
	}
	return &operationImpl{
		id: OP12,
		f:  createOps6_12(setup, complexAssemblyOperation),
	}
}

func createOperation13(setup Setup) Operation {
	baseAssemblyOperation := func(tx Transaction, assembly BaseAssembly) {
		assembly.UpdateBuildDate(tx)
	}
	return &operationImpl{
		id: OP13,
		f:  createOps7_13(setup, baseAssemblyOperation),
	}
}

func createOperation14(setup Setup) Operation {
	componentOperation := func(tx Transaction, component CompositePart) {
		component.UpdateBuildDate(tx)
	}
	return &operationImpl{
		id: OP14,
		f:  createOps8_14(setup, componentOperation),
	}
}

func createOperation15(setup Setup) Operation {
	atomicPartOperation := func(tx Transaction, part AtomicPart) {
		RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
		part.UpdateBuildDate(tx)
		AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
	}
	return &operationImpl{
		id: OP15,
		f:  createOps1_9_15(setup, atomicPartOperation),
	}
}

func createOps1_9_15(setup Setup, atomicPartOperation func(tx Transaction, part AtomicPart)) func(tx Transaction) (int, error) {
	f := func(tx Transaction) (int, error) {
		count := 0
		for i := 0; i < 10; i++ {
			partID := rand.Intn(internal.MaxAtomicParts) + 1
			part := setup.AtomicPartIDIndex().Get(tx, partID)
			if part == nil {
				continue
			}
			atomicPartOperation(tx, part.(AtomicPart))
			count++
		}
		return count, nil
	}
	return f
}

func createOps2_3_10(setup Setup, percent int, atomicPartOperation func(tx Transaction, part AtomicPart)) func(tx Transaction) (int, error) {
	minAtomicDate := internal.MaxAtomicPartDate - percent*(internal.MaxAtomicPartDate-internal.MinAtomicPartDate)/100
	f := func(tx Transaction) (int, error) {
		partSets := setup.AtomicPartBuildDateIndex().GetRange(tx, minAtomicDate, internal.MaxAtomicPartDate)
		count := 0
		for _, partSet := range partSets {
			set := partSet.(LargeSet)
			for _, part := range set.ToSlice(tx) {
				atomicPartOperation(tx, part.(AtomicPart))
				count++
			}
		}
		return count, nil
	}
	return f
}

func createOps6_12(setup Setup, complexAssemblyOperation func(tx Transaction, assembly ComplexAssembly)) func(tx Transaction) (int, error) {
	f := func(tx Transaction) (int, error) {
		complexAssemblyID := rand.Intn(internal.MaxComplexAssemblies) + 1
		complexAssembly := setup.ComplexAssemblyIDIndex().Get(tx, complexAssemblyID)
		if complexAssembly == nil {
			return -1, NewOpFailedError("")
		}

		superAssembly := complexAssembly.(ComplexAssembly).GetSuperAssembly(tx)
		if superAssembly == nil {
			complexAssemblyOperation(tx, complexAssembly.(ComplexAssembly))
			return 1, nil
		}

		count := 0
		for _, siblingAssembly := range superAssembly.GetSubAssemblies(tx).ToSlice() {
			complexAssemblyOperation(tx, siblingAssembly.(ComplexAssembly))
			count++
		}
		return count, nil
	}
	return f
}

func createOps7_13(setup Setup, baseAssemblyOperation func(tx Transaction, assembly BaseAssembly)) func(tx Transaction) (int, error) {
	f := func(tx Transaction) (int, error) {
		baseAssemblyID := rand.Intn(internal.MaxBaseAssemblies) + 1
		baseAssembly := setup.BaseAssemblyIDIndex().Get(tx, baseAssemblyID)
		if baseAssembly == nil {
			return -1, NewOpFailedError("")
		}

		superAssembly := baseAssembly.(BaseAssembly).GetSuperAssembly(tx)

		count := 0
		for _, siblingAssembly := range superAssembly.GetSubAssemblies(tx).ToSlice() {
			baseAssemblyOperation(tx, siblingAssembly.(BaseAssembly))
			count++
		}
		return count, nil
	}
	return f
}

func createOps8_14(setup Setup, componentOperation func(tx Transaction, component CompositePart)) func(tx Transaction) (int, error) {
	f := func(tx Transaction) (int, error) {
		baseAssemblyID := rand.Intn(internal.MaxBaseAssemblies) + 1
		baseAssembly := setup.BaseAssemblyIDIndex().Get(tx, baseAssemblyID)
		if baseAssembly == nil {
			return -1, NewOpFailedError("")
		}

		count := 0
		for _, component := range baseAssembly.(BaseAssembly).GetComponents(tx).ToSlice() {
			componentOperation(tx, component.(CompositePart))
			count++
		}
		return count, nil
	}
	return f
}
