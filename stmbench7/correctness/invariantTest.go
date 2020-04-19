package correctness

import (
    "fmt"
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
    "os"
)

type traversedObjects struct {
    complexAssemblies map[interface{}]struct{}
    baseAssemblies map[interface{}]struct{}
    components map[interface{}]struct{}
    documents map[interface{}]struct{}
    atomicParts map[interface{}]struct{}
}

func CheckInvariants(tx Transaction, setup Setup, initial bool) error {
    objects := traversedObjects{
        complexAssemblies: map[interface{}]struct{}{},
        baseAssemblies: map[interface{}]struct{}{},
        components: map[interface{}]struct{}{},
        documents: map[interface{}]struct{}{},
        atomicParts: map[interface{}]struct{}{},
    }

    fmt.Fprintln(os.Stderr, "Checking data structures...")
    if err := moduleInvariantTest(tx, setup.Module(), initial, objects); err != nil {
        return newInvFailedError(setup.Module(), setup.Module().GetId(tx), err.Error())
    }

    fmt.Fprintln(os.Stderr, "Checking indexes...")
    if err := indexInvariantTest(tx, setup, initial, objects); err != nil {
        return err
    }
    fmt.Fprintln(os.Stderr, "Invariants ok.")
    fmt.Fprintln(os.Stdout, "Invariants ok.")
    return nil
}

func indexInvariantTest(tx Transaction, setup Setup, initial bool, objects traversedObjects) error {
    // Complex assemblies
    complexAssemblyIndex := setup.ComplexAssemblyIDIndex()
    for traversedAssembly := range objects.complexAssemblies {
        traversed := traversedAssembly.(ComplexAssembly)
        id := traversed.GetId(tx)
        indexed := complexAssemblyIndex.Get(tx, id)
        if err := checkIndexValue("ComplexAssembly.ID", traversed, indexed, id); err != nil {
            return err
        }
    }
    if err := checkAllTraversed(tx, complexAssemblyIndex, objects.complexAssemblies, "ComplexAssembly.ID"); err != nil {
        return err
    }

    // Base assemblies
    baseAssemblyIndex := setup.BaseAssemblyIDIndex()
    for traversedAssembly := range objects.baseAssemblies {
        traversed := traversedAssembly.(BaseAssembly)
        id := traversed.GetId(tx)
        indexed := baseAssemblyIndex.Get(tx, id)
        if err := checkIndexValue("BaseAssembly.ID", traversed, indexed, id); err != nil {
            return err
        }
    }
    if err := checkAllTraversed(tx, baseAssemblyIndex, objects.baseAssemblies, "BaseAssembly.ID"); err != nil {
        return err
    }

    // Composite parts
    compositePartIndex := setup.CompositePartIDIndex()
    for traversedComponent := range objects.components {
        traversed := traversedComponent.(CompositePart)
        id := traversed.GetId(tx)
        indexed := compositePartIndex.Get(tx, id)
        if err := checkIndexValue("CompositePart.ID", traversed, indexed, id); err != nil {
            return err
        }
    }
    // Check invariants for components disconnected from the data structure
    // (and add those components to the set of traversed objects)
    for _, key := range compositePartIndex.GetKeys(tx) {
        if _, ok := objects.components[compositePartIndex.Get(tx, key)]; !ok {
            if err := compositePartInvariantTest(tx, compositePartIndex.Get(tx, key).(CompositePart), initial, nil, objects); err != nil {
                return err
            }
        }
    }

    // Documents
    documentTitleIndex := setup.DocumentTitleIndex()
    for traversedDocument := range objects.documents {
        traversed := traversedDocument.(Document)
        title := traversed.GetTitle(tx)
        id := traversed.GetDocumentId(tx)
        indexed := documentTitleIndex.Get(tx, title)
        if err := checkIndexValue("Document.title", traversed, indexed, id); err != nil {
            return err
        }
    }
    if err := checkAllTraversed(tx, documentTitleIndex, objects.documents, "Document.title"); err != nil {
        return err
    }

    // Atomic parts (id index)
    atomicPartIDIndex := setup.AtomicPartIDIndex()
    for traversedPart := range objects.atomicParts {
        traversed := traversedPart.(AtomicPart)
        id := traversed.GetId(tx)
        indexed := atomicPartIDIndex.Get(tx, id)
        if err := checkIndexValue("AtomicPart.ID", traversed, indexed, id); err != nil {
            return err
        }
    }
    if err := checkAllTraversed(tx, atomicPartIDIndex, objects.atomicParts, "AtomicPart.ID"); err != nil {
        return err
    }

    // Atomic parts (buildDate index)
    atomicPartBuildDateIndex := setup.AtomicPartBuildDateIndex()
    for traversedPart := range objects.atomicParts {
        traversed := traversedPart.(AtomicPart)
        id := traversed.GetId(tx)
        sameBuildDateParts := atomicPartBuildDateIndex.Get(tx, traversed.GetBuildDate(tx)).(LargeSet)
        if sameBuildDateParts == nil || !sameBuildDateParts.Contains(tx, traversed) {
            return newIndexError("AtomicPart.buildDate", fmt.Sprintf("element with id %d not in the index", id))
        }
    }
    for _, key := range atomicPartBuildDateIndex.GetKeys(tx) {
        for _, indexed := range atomicPartBuildDateIndex.Get(tx, key).(LargeSet).ToSlice(tx) {
            if _, ok := objects.atomicParts[indexed]; !ok {
                return newIndexError("AtomicPart.buildDate", "index contains too many elements")
            }
        }
    }

    return nil
}