package operations

import (
    "fmt"
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
    "math/rand"
)

func createShortTraversal1(setup Setup) Operation {
    componentTraverse := func(tx Transaction, component CompositePart) int {
        parts := component.GetParts(tx).ToSlice(tx)
        numOfAtomicParts := len(parts)
        if numOfAtomicParts == 0 {
            panic("ST1: illegal size of CompositePart.parts!")
        }
        part := parts[rand.Intn(numOfAtomicParts)].(AtomicPart)
        return part.GetX(tx) + part.GetY(tx)
    }
    return &operationImpl{
        id: ST1,
        f:  createShortTraversals1_2_6_7_9_10(setup, componentTraverse),
    }
}

func createShortTraversal2(setup Setup) Operation {
    componentTraverse := func(tx Transaction, component CompositePart) int {
        return component.GetDocumentation(tx).SearchText(tx, 'I')
    }
    return &operationImpl{
        id: ST2,
        f:  createShortTraversals1_2_6_7_9_10(setup, componentTraverse),
    }
}

func createShortTraversal3(setup Setup) Operation {
    assemblyOperation := func(tx Transaction, assembly Assembly) {
        assembly.NullOperation(tx)
    }
    return &operationImpl{
        id: ST3,
        f:  createShortTraversals3_8(setup, assemblyOperation),
    }
}

func createShortTraversal4(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        result := 0

        for i := 0 ; i < 100 ; i++ {
            partID := rand.Intn(internal.MaxCompParts) + 1
            docTitle := fmt.Sprintf("Composite Part #%d", partID)

            document := setup.DocumentTitleIndex().Get(tx, docTitle)
            if document == nil {
                continue
            }

            for _, assembly := range document.(Document).GetCompositePart(tx).GetUsedIn(tx).ToSlice() {
                assembly.(Assembly).NullOperation(tx)
                result++
            }
        }

        return result, nil
    }
    return &operationImpl{
        id: ST4,
        f:  f,
    }
}

func createShortTraversal5(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        result := 0
        for _, key := range setup.BaseAssemblyIDIndex().GetKeys(tx) {
            result += checkBaseAssembly(tx, setup.BaseAssemblyIDIndex().Get(tx, key).(BaseAssembly))
        }
        return result, nil
    }
    return &operationImpl{
        id: ST5,
        f:  f,
    }
}

func createShortTraversal6(setup Setup) Operation {
    componentTraverse := func(tx Transaction, component CompositePart) int {
        parts := component.GetParts(tx).ToSlice(tx)
        numOfAtomicParts := len(parts)
        if numOfAtomicParts == 0 {
            panic("ST1: illegal size of CompositePart.parts!")
        }
        part := parts[rand.Intn(numOfAtomicParts)].(AtomicPart)
        part.SwapXY(tx)
        return part.GetX(tx) + part.GetY(tx)
    }
    return &operationImpl{
        id: ST6,
        f:  createShortTraversals1_2_6_7_9_10(setup, componentTraverse),
    }
}

func createShortTraversal7(setup Setup) Operation {
    componentTraverse := func(tx Transaction, component CompositePart) int {
        document := component.GetDocumentation(tx)
        if document.TextBeginsWith(tx, "I am") {
            return document.ReplaceText(tx, "I am", "This is")
        } else if document.TextBeginsWith(tx, "This is") {
            return document.ReplaceText(tx, "This is", "I am")
        }
        panic("ST7: unexpected beginning of Document.text!")
    }
    return &operationImpl{
        id: ST7,
        f:  createShortTraversals1_2_6_7_9_10(setup, componentTraverse),
    }
}

func createShortTraversal8(setup Setup) Operation {
    assemblyOperation := func(tx Transaction, assembly Assembly) {
        assembly.UpdateBuildDate(tx)
    }
    return &operationImpl{
        id: ST8,
        f:  createShortTraversals3_8(setup, assemblyOperation),
    }
}

func createShortTraversal9(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        part.SwapXY(tx)
        return 1
    }
    componentTraverse := func(tx Transaction, component CompositePart) int {
        return t1TraverseAtomicPart(tx, component.GetRootPart(tx), map[AtomicPart]struct{}{}, atomicPartOperation)
    }
    return &operationImpl{
        id: ST9,
        f:  createShortTraversals1_2_6_7_9_10(setup, componentTraverse),
    }
}

func createShortTraversal10(setup Setup) Operation {
    atomicPartOperation := func(tx Transaction, part AtomicPart, visited map[AtomicPart]struct{}) int {
        part.NullOperation(tx)
        return 1
    }
    componentTraverse := func(tx Transaction, component CompositePart) int {
        return t1TraverseAtomicPart(tx, component.GetRootPart(tx), map[AtomicPart]struct{}{}, atomicPartOperation)
    }
    return &operationImpl{
        id: ST10,
        f:  createShortTraversals1_2_6_7_9_10(setup, componentTraverse),
    }
}

func createShortTraversals1_2_6_7_9_10(setup Setup, componentTraverse func(tx Transaction, component CompositePart) int) func(tx Transaction) (int, error) {
    f := func(tx Transaction) (int, error) {
        return randomTraverseAssembly(tx, setup.Module().GetDesignRoot(tx), componentTraverse)
    }
    return f
}

func randomTraverseAssembly(tx Transaction, assembly Assembly, componentTraverse func(tx Transaction, component CompositePart) int) (int, error) {
    switch a := assembly.(type) {
    case BaseAssembly:
        return randomTraverseBaseAssembly(tx, a, componentTraverse)
    case ComplexAssembly:
        return randomTraverseComplexAssembly(tx, a, componentTraverse)
    }
    panic("Assembly does not have valid type")
}

func randomTraverseComplexAssembly(tx Transaction, assembly ComplexAssembly, componentTraverse func(tx Transaction, component CompositePart) int) (int, error) {
    subAssemblies := assembly.GetSubAssemblies(tx).ToSlice()
    nextAssembly := rand.Intn(len(subAssemblies))
    return randomTraverseAssembly(tx, subAssemblies[nextAssembly].(Assembly), componentTraverse)
}

func randomTraverseBaseAssembly(tx Transaction, assembly BaseAssembly, componentTraverse func(tx Transaction, component CompositePart) int) (int, error) {
    components := assembly.GetComponents(tx).ToSlice()
    numOfComponents := len(components)
    if numOfComponents == 0 {
        return 0, NewOpFailedError("")
    }
    nextComponent := rand.Intn(numOfComponents)
    return componentTraverse(tx, components[nextComponent].(CompositePart)), nil
}

func createShortTraversals3_8(setup Setup, assemblyOperation func(tx Transaction, assembly Assembly)) func(tx Transaction) (int, error) {
    f := func(tx Transaction) (int, error) {
        partID := rand.Intn(internal.MaxAtomicParts) + 1
        part := setup.AtomicPartIDIndex().Get(tx, partID)
        if part == nil {
            return 0, nil
        }
        return st3TraverseComponent(tx, part.(AtomicPart).GetPartOf(tx), assemblyOperation), nil
    }
    return f
}

func st3TraverseComponent(tx Transaction, part CompositePart, assemblyOperation func(tx Transaction, assembly Assembly)) int {
    visited := map[Assembly]struct{}{}
    result := 0
    for _, assembly := range part.GetUsedIn(tx).ToSlice() {
        result += st3TraverseAssembly(tx, assembly.(Assembly), visited, assemblyOperation)
    }
    return result
}

func st3TraverseAssembly(tx Transaction, assembly Assembly, visited map[Assembly]struct{}, assemblyOperation func(tx Transaction, assembly Assembly)) int {
    if assembly == nil {
        return 0
    }
    if _, ok := visited[assembly]; ok {
        return 0
    }

    visited[assembly] = struct{}{}
    assemblyOperation(tx, assembly)
    return st3TraverseAssembly(tx, assembly.GetSuperAssembly(tx), visited, assemblyOperation)
}
