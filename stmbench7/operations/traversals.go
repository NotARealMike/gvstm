package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

func createTraversal1(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        part.NullOperation(tx)
        return 1
    }
    return &operationImpl{
        id: T1,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal2a(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        if len(visited) == 0 {
            part.SwapXY(tx)
            return 1
        }
        part.NullOperation(tx)
        return 0
    }
    return &operationImpl{
        id: T2a,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal2b(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        part.SwapXY(tx)
        return 1
    }
    return &operationImpl{
        id: T2b,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal2c(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        part.SwapXY(tx)
        part.SwapXY(tx)
        part.SwapXY(tx)
        part.SwapXY(tx)
        return 4
    }
    return &operationImpl{
        id: T2c,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal3a(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        if len(visited) == 0 {
            RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
            part.UpdateBuildDate(tx)
            AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
            return 1
        }
        part.NullOperation(tx)
        return 0
    }
    return &operationImpl{
        id: T3a,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal3b(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        part.UpdateBuildDate(tx)
        AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        return 1
    }
    return &operationImpl{
        id: T3b,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal3c(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        part.UpdateBuildDate(tx)
        AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        part.UpdateBuildDate(tx)
        AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        part.UpdateBuildDate(tx)
        AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        RemoveAtomicPartFromBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        part.UpdateBuildDate(tx)
        AddAtomicPartToBuildDateIndex(tx, setup.AtomicPartBuildDateIndex(), part)
        return 4
    }
    return &operationImpl{
        id: T3c,
        f:  createTraversals1_2_3(setup, atomicPartOperation),
    }
}

func createTraversal4(setup Setup) Operation {
    documentOperation := func(tx Transaction, document Document) int {
        return document.SearchText(tx, 'I')
    }
    return &operationImpl{
        id: T4,
        f:  createTraversals4_5(setup, documentOperation),
    }
}

func createTraversal5(setup Setup) Operation {
    documentOperation := func(tx Transaction, document Document) int {
        if document.TextBeginsWith(tx, "I am") {
            return document.ReplaceText(tx, "I am", "This is")
        } else if document.TextBeginsWith(tx,"This is") {
            return document.ReplaceText(tx, "This is", "I am")
        }
        panic("T5: Illegal document text: " + document.GetText(tx))
    }
    return &operationImpl{
        id: T5,
        f:  createTraversals4_5(setup, documentOperation),
    }
}

func createTraversal6(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        return t6TraverseComplexAssembly(tx, setup.Module().GetDesignRoot(tx)), nil
    }
    return &operationImpl{
        id: T6,
        f:  f,
    }
}

func createTraversalQ6(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        return checkComplexAssembly(tx, setup.Module().GetDesignRoot(tx)), nil
    }
    return &operationImpl{
        id: Q6,
        f:  f,
    }
}

func createTraversalQ7(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        result := 0

        for _, key := range setup.AtomicPartIDIndex().GetKeys(tx) {
            setup.AtomicPartIDIndex().Get(tx, key).(AtomicPart).NullOperation(tx)
            result++
        }

        return result, nil
    }
    return &operationImpl{
        id: Q7,
        f:  f,
    }
}

func createTraversals1_2_3(setup Setup, atomicPartOperation func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int) func(tx Transaction) (int, error) {
    f := func(tx Transaction) (int, error) {
        return t1TraverseComplexAssembly(tx, setup.Module().GetDesignRoot(tx), atomicPartOperation), nil
    }
    return f
}

func t1TraverseAssembly(tx Transaction, assembly Assembly, atomicPartOperation func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int) int {
    switch a := assembly.(type) {
    case BaseAssembly:
        return t1TraverseBaseAssembly(tx, a, atomicPartOperation)
    case ComplexAssembly:
        return t1TraverseComplexAssembly(tx, a, atomicPartOperation)
    }
    panic("Assembly does not have valid type")
}

func t1TraverseComplexAssembly(tx Transaction, assembly ComplexAssembly, atomicPartOperation func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int) int {
    partsVisited := 0
    for _, subAssemby := range assembly.GetSubAssemblies(tx).ToSlice() {
        partsVisited += t1TraverseAssembly(tx, subAssemby.(Assembly), atomicPartOperation)
    }
    return partsVisited
}

func t1TraverseBaseAssembly(tx Transaction, assembly BaseAssembly, atomicPartOperation func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int) int {
    partsVisited := 0
    for _, component := range assembly.GetComponents(tx).ToSlice() {
        partsVisited += t1TraverseCompositePart(tx, component.(CompositePart), atomicPartOperation)
    }
    return partsVisited
}

func t1TraverseCompositePart(tx Transaction, component CompositePart, atomicPartOperation func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int) int {
    return t1TraverseAtomicPart(tx, component.GetRootPart(tx), map[AtomicPart]struct{}{}, atomicPartOperation)
}

func t1TraverseAtomicPart(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}, atomicPartOperation func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int) int {
    if part == nil {
        return 0
    }
    if _, ok := visited[part]; ok {
        return 0
    }

    result := atomicPartOperation(tx, part, visited)
    visited[part] = struct{}{}

    for _, connection := range part.GetToConnections(tx).ToSlice() {
        result += t1TraverseAtomicPart(tx, connection.(Connection).GetDestination(), visited, atomicPartOperation)
    }
    return result
}

func createTraversals4_5(setup Setup, documentOperation func(tx Transaction, document Document) int) func(tx Transaction) (int, error) {
    f := func(tx Transaction) (int, error) {
        return t4TraverseComplexAssembly(tx, setup.Module().GetDesignRoot(tx), documentOperation), nil
    }
    return f
}

func t4TraverseAssembly(tx Transaction, assembly Assembly, documentOperation func(tx Transaction, document Document) int) int {
    switch a := assembly.(type) {
    case BaseAssembly:
        return t4TraverseBaseAssembly(tx, a, documentOperation)
    case ComplexAssembly:
        return t4TraverseComplexAssembly(tx, a, documentOperation)
    }
    panic("Assembly does not have valid type")

}

func t4TraverseComplexAssembly(tx Transaction, assembly ComplexAssembly, documentOperation func(tx Transaction, document Document) int) int {
    count := 0
    for _, subAssemby := range assembly.GetSubAssemblies(tx).ToSlice() {
        count += t4TraverseAssembly(tx, subAssemby.(Assembly), documentOperation)
    }
    return count
}

func t4TraverseBaseAssembly(tx Transaction, assembly BaseAssembly, documentOperation func(tx Transaction, document Document) int) int {
    count := 0
    for _, component := range assembly.GetComponents(tx).ToSlice() {
        count += documentOperation(tx, component.(CompositePart).GetDocumentation(tx))
    }
    return count
}

func t6TraverseAssembly(tx Transaction, assembly Assembly) int {
    switch a := assembly.(type) {
    case BaseAssembly:
        return t6TraverseBaseAssembly(tx, a)
    case ComplexAssembly:
        return t6TraverseComplexAssembly(tx, a)
    }
    panic("Assembly does not have valid type")
}

func t6TraverseComplexAssembly(tx Transaction, assembly ComplexAssembly) int {
    partsVisited := 0
    for _, subAssemby := range assembly.GetSubAssemblies(tx).ToSlice() {
        partsVisited += t6TraverseAssembly(tx, subAssemby.(Assembly))
    }
    return partsVisited
}

func t6TraverseBaseAssembly(tx Transaction, assembly BaseAssembly) int {
    partsVisited := 0
    for _, component := range assembly.GetComponents(tx).ToSlice() {
        component.(CompositePart).GetRootPart(tx).NullOperation(tx)
        partsVisited += 1
    }
    return partsVisited
}

func checkAssembly(tx Transaction, assembly Assembly) int {
    switch a := assembly.(type) {
    case BaseAssembly:
        return checkBaseAssembly(tx, a)
    case ComplexAssembly:
        return checkComplexAssembly(tx, a)
    }
    panic("Assembly does not have valid type")
}

func checkBaseAssembly(tx Transaction, assembly BaseAssembly) int {
    assBuildDate := assembly.GetBuildDate(tx)

    for _, part := range assembly.GetComponents(tx).ToSlice() {
        if part.(CompositePart).GetBuildDate(tx) > assBuildDate {
            assembly.NullOperation(tx)
            return 1
        }
    }
    return 0
}

func checkComplexAssembly(tx Transaction, assembly ComplexAssembly) int {
    result := 0

    for _, subAssembly := range assembly.GetSubAssemblies(tx).ToSlice() {
        result += checkAssembly(tx, subAssembly.(Assembly))
    }

    if result == 0 {
        return 0
    }

    assembly.NullOperation(tx)
    return result + 1
}