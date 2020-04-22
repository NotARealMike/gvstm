package operations

import (
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
    "gvstm/stmbench7/internal"
    "math/rand"
)

func createStructureModification1(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        _, err := setup.CompositePartBuilder().CreateAndRegisterCompositePart(tx)
        return 0, err
    }
    return &operationImpl{
        id: SM1,
        f:  f,
    }
}

func createStructureModification2(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        partID := rand.Intn(internal.MaxCompParts) + 1
        partToRemove := setup.CompositePartIDIndex().Get(tx, partID)
        if partToRemove == nil {
            return 0, NewOpFailedError("")
        }
        setup.CompositePartBuilder().UnregisterAndRecycleCompositePart(tx, partToRemove.(CompositePart))
        return 0, nil
    }
    return &operationImpl{
        id: SM2,
        f:  f,
    }
}

func createStructureModification3(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        baseAssemblyID := rand.Intn(internal.MaxBaseAssemblies) + 1
        componentID := rand.Intn(internal.MaxCompParts) + 1
        baseAssembly := setup.BaseAssemblyIDIndex().Get(tx, baseAssemblyID)
        component := setup.CompositePartIDIndex().Get(tx, componentID)

        if baseAssembly == nil || component == nil {
            return 0, NewOpFailedError("")
        }

        baseAssembly.(BaseAssembly).AddComponent(tx, component.(CompositePart))

        return 0, nil
    }
    return &operationImpl{
        id: SM3,
        f:  f,
    }
}

func createStructureModification4(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        baseAssemblyID := rand.Intn(internal.MaxBaseAssemblies) + 1
        baseAssembly := setup.BaseAssemblyIDIndex().Get(tx, baseAssemblyID)
        if baseAssembly == nil {
            return 0, NewOpFailedError("")
        }

        components := baseAssembly.(BaseAssembly).GetComponents(tx)
        numOfComponents := components.Size()
        if numOfComponents == 0 {
            return 0, NewOpFailedError("")
        }

        componentToRemove := rand.Intn(numOfComponents)
        baseAssembly.(BaseAssembly).RemoveComponent(tx, components.ToSlice()[componentToRemove].(CompositePart))
        return 0, nil
    }
    return &operationImpl{
        id: SM4,
        f:  f,
    }
}

func createStructureModification5(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        siblingAssemblyID := rand.Intn(internal.MaxBaseAssemblies) + 1
        siblingAssembly := setup.BaseAssemblyIDIndex().Get(tx, siblingAssemblyID)
        if siblingAssembly == nil {
            return 0, NewOpFailedError("")
        }

        superAssembly := siblingAssembly.(BaseAssembly).GetSuperAssembly(tx)
        _, err := setup.AssemblyBuilder().CreateAndRegisterAssembly(tx, setup.Module(), superAssembly)
        return 0, err
    }
    return &operationImpl{
        id: SM5,
        f:  f,
    }
}

func createStructureModification6(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        baseAssemblyID := rand.Intn(internal.MaxBaseAssemblies) + 1
        baseAssembly := setup.BaseAssemblyIDIndex().Get(tx, baseAssemblyID)
        if baseAssembly == nil {
            return 0, NewOpFailedError("")
        }

        if baseAssembly.(BaseAssembly).GetSuperAssembly(tx).GetSubAssemblies(tx).Size() == 1 {
            return 0, NewOpFailedError("")
        }
        setup.AssemblyBuilder().UnregisterAndRecycleAssembly(tx, baseAssembly.(BaseAssembly))
        return 1, nil
    }
    return &operationImpl{
        id: SM6,
        f:  f,
    }
}

func createStructureModification7(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        superAssemblyID := rand.Intn(internal.MaxComplexAssemblies) + 1
        superAssembly := setup.ComplexAssemblyIDIndex().Get(tx, superAssemblyID)
        if superAssembly == nil {
            return 0, NewOpFailedError("")
        }
        _, err := setup.AssemblyBuilder().CreateAndRegisterAssembly(tx, setup.Module(), superAssembly.(ComplexAssembly))
        return 1, err
    }
    return &operationImpl{
        id: SM7,
        f:  f,
    }
}

func createStructureModification8(setup Setup) Operation {
    f := func(tx Transaction) (int, error) {
        complexAssemblyID := rand.Intn(internal.MaxComplexAssemblies) + 1
        complexAssembly := setup.ComplexAssemblyIDIndex().Get(tx, complexAssemblyID)
        if complexAssembly == nil {
            return 0, NewOpFailedError("")
        }

        superAssembly := complexAssembly.(ComplexAssembly).GetSuperAssembly(tx)
        if superAssembly == nil || superAssembly.GetSubAssemblies(tx).Size() == 1 {
            return 0, NewOpFailedError("")
        }

        setup.AssemblyBuilder().UnregisterAndRecycleAssembly(tx, complexAssembly.(ComplexAssembly))
        return 1, nil
    }
    return &operationImpl{
        id: SM8,
        f:  f,
    }
}
