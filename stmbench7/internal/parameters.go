package internal

import "math"

var (
    NumAtomsPerCompPart = 200
    NumConnectionsPerAtomicPart = 6
    DocumentSize = 20000
    ManualSize = 1000000
    NumCompPartsPerModule = 50//500
    NumSubAssemblies  = 3
    NumAssemblyLevels = 7
    NumCompPartPerAss = 3
    NumModules = 1

    InitialTotalCompParts = NumModules * NumCompPartsPerModule
    InitialTotalBaseAssemblies = int(math.Pow(float64(NumSubAssemblies), float64(NumAssemblyLevels-1)))
    InitialTotalComplexAssemblies = (1 - InitialTotalBaseAssemblies) / (1 - NumSubAssemblies)

    MaxCompParts = int(1.05 * float64(InitialTotalCompParts))
    MaxAtomicParts = MaxCompParts * NumAtomsPerCompPart
    MaxBaseAssemblies = int(1.05 * float64(InitialTotalBaseAssemblies))
    MaxComplexAssemblies = int(1.05 * float64(InitialTotalComplexAssemblies))

    MinModuleDate = 1000
    MaxModuleDate = 1999
    MinAssDate = 1000
    MaxAssDate = 1999
    MinAtomicPartDate = 1000
    MaxAtomicPartDate = 1999
    MinOldCompositePartDate = 0
    MaxOldCompositePartDate = 999
    MinYoungCompositePartDate = 2000
    MaxYoungCompositePartDate = 2999
    YoungCompositePartFraction = 10
    NumTypes = 10
    XYRange = 100000

    TraversalsRatio = 5
    ShortTraversalsRatio = 40
    OperationsRatio = 45
    StructuralModificationsRatio = 10

    ReadOnlyWorkloadRORatio = 100
    ReadDominatedWorkloadRORatio = 90
    ReadWriteWorkloadRORatio = 60
    WriteDominatedWorkloadRORatio = 10

    MaxLowTTC = 999
    HighTTCEntries = 200
    HighTTCLogBase = 1.03
)
