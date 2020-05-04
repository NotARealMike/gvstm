package main

import (
	"fmt"
	"gvstm/gvstm"
	"gvstm/stm"
	"gvstm/stmbench7/correctness"
	"gvstm/stmbench7/interfaces"
	"gvstm/stmbench7/internal"
	"gvstm/stmbench7/operations"
	"math"
	"os"
	"sync"
	"time"
)

type benchmark interface {
	createInitialClone()
	start()
	checkInvariants(initial bool) error
	checkOpacity() error
	showTTCHistograms()
	showStats()
}

type benchmarkImpl struct {
	params *benchmarkParams

	elapsedTime  time.Duration
	operationCDF []float64
	benchThreads []benchThread
	setup        interfaces.Setup
}

func createBenchmark(params *benchmarkParams) benchmark {
	b := &benchmarkImpl{params: params}
	interfaces.SetFactories(b.params.initialiser)
	operations.OEFactory = params.executorFactory
	fmt.Fprintln(os.Stderr, header())
	printOutAndErr(runtimeParamsInfo(b.params))
	b.generateOperationCDF()
	b.setupStructures()
	return b
}

func (b *benchmarkImpl) setupStructures() {
	fmt.Fprintln(os.Stderr, "Setup start...")
	if b.params.gvstm {
		gvstm.Atomic(func(tx stm.Transaction) {
			b.setup = interfaces.NewSetup(tx)
		})
	} else {
		b.setup = interfaces.NewSetup(nil)
	}
	b.benchThreads = make([]benchThread, b.params.numThreads)
	for i := 0; i < b.params.numThreads; i++ {
		b.benchThreads[i] = createBenchThread(b.setup, b.operationCDF, i)
	}
	fmt.Fprintln(os.Stderr, "Setup completed.")
}

func (b *benchmarkImpl) createInitialClone() {}

func (b *benchmarkImpl) start() {
	fmt.Fprintln(os.Stderr, "\nBenchmark started...")
	startTime := time.Now()
	var wg sync.WaitGroup
	for _, t := range b.benchThreads {
		wg.Add(1)
		go t.run(&wg)
	}
	time.Sleep(b.params.duration)
	for _, t := range b.benchThreads {
		t.stopThread()
	}
	wg.Wait()
	b.elapsedTime = time.Now().Sub(startTime)
	fmt.Fprintln(os.Stderr, "Benchmark completed.\n")
}

func (b *benchmarkImpl) generateOperationCDF() {
	stRatio := float64(internal.ShortTraversalsRatio) / 100
	opRatio := float64(internal.OperationsRatio) / 100
	tvRatio, smRatio := 0., 0.
	if b.params.traversalsEnabled {
		tvRatio = float64(internal.TraversalsRatio) / 100
	}
	if b.params.structureModificationsEnabled {
		smRatio = float64(internal.StructuralModificationsRatio) / 100
	}

	readOnlyRatio := float64(b.params.readOnlyRatio) / 100
	updateRatio := 1 - readOnlyRatio

	roTraversalsRatio := tvRatio * readOnlyRatio
	rwTraversalsRatio := 0.
	if b.params.longReadWriteTraversalsEnabled {
		rwTraversalsRatio = tvRatio * updateRatio
	}

	sumRatios := roTraversalsRatio + rwTraversalsRatio + stRatio + opRatio + smRatio*updateRatio

	roTraversalsRatio /= sumRatios
	rwTraversalsRatio /= sumRatios
	stRatio /= sumRatios
	opRatio /= sumRatios
	smRatio /= sumRatios

	for _, opID := range operations.OperationIDs {
		opID.OpType.Count++
	}

	operations.TraversalRW.Probability = rwTraversalsRatio / float64(operations.TraversalRW.Count)
	operations.TraversalRO.Probability = roTraversalsRatio / float64(operations.TraversalRO.Count)
	operations.ShortTraversalRW.Probability = stRatio * updateRatio / float64(operations.ShortTraversalRW.Count)
	operations.ShortTraversalRO.Probability = stRatio * readOnlyRatio / float64(operations.ShortTraversalRO.Count)
	operations.OperationRW.Probability = opRatio * updateRatio / float64(operations.OperationRW.Count)
	operations.OperationRO.Probability = opRatio * readOnlyRatio / float64(operations.OperationRO.Count)
	operations.StructureModification.Probability = smRatio * updateRatio / float64(operations.StructureModification.Count)

	printOutAndErr("Operation ratios [%]:\n")
	for _, opType := range operations.OperationTypes {
		printOutAndErr(alignText(opType.Name, 23) + ": " + alignText(formatFloat(opType.Probability*float64(opType.Count)*100), 6) + "\n")
	}
	printOutAndErr("\n")

	operationProbabilities := make([]float64, len(operations.OperationIDs))
	for i := range operations.OperationIDs {
		operationProbabilities[i] = operations.OperationIDs[i].OpType.Probability
	}

	b.operationCDF = make([]float64, len(operations.OperationIDs))
	b.operationCDF[0] = operationProbabilities[0]
	for i := 1; i < len(operationProbabilities); i++ {
		b.operationCDF[i] = b.operationCDF[i-1] + operationProbabilities[i]
	}
	b.operationCDF[len(operationProbabilities)-1] = 1
}

func (b *benchmarkImpl) checkInvariants(initial bool) (err error) {
	if initial {
		printOutAndErr("Checking invariants (initial data structure):\n")
	} else {
		printOutAndErr("Checking invariants (final data structure):\n")
	}
	if b.params.gvstm {
		gvstm.Atomic(func(tx stm.Transaction) {
			err = correctness.CheckInvariants(tx, b.setup, initial)
		})
	} else {
		err = correctness.CheckInvariants(nil, b.setup, initial)
	}
	if err == nil {
		printOutAndErr("Invariant check completed successfully!\n")
	}
	return
}

func (b *benchmarkImpl) checkOpacity() error {
	return nil
}

func (b *benchmarkImpl) showTTCHistograms() {
	if !b.params.printTTCHistograms {
		return
	}
	printOutAndErr(section("TTC histograms"))

	for i := range operations.OperationIDs {
		printOutAndErr(fmt.Sprintf("TTC histogram for %s :", operations.OperationIDs[i].Name))

		for ttc := 0; ttc <= internal.MaxLowTTC; ttc++ {
			count := 0
			for _, thread := range b.benchThreads {
				count += thread.opsTTC(i, ttc)
			}
			printOutAndErr(fmt.Sprintf(" %d,%d", ttc, count))
		}

		for logTTCIndex := 0; logTTCIndex < internal.HighTTCEntries; logTTCIndex++ {
			count := 0
			for _, thread := range b.benchThreads {
				count += thread.opsHighTTCLog(i, logTTCIndex)
			}
			ttc := logTTCIndexToTTC(logTTCIndex)
			printOutAndErr(fmt.Sprintf(" %d,%d", ttc, count))
		}
		printOutAndErr("\n")
	}
	printOutAndErr("\n")
}

func (b *benchmarkImpl) showStats() {
	printOutAndErr(section("Detailed results"))

	for opIDIndex, opID := range operations.OperationIDs {
		printOutAndErr(fmt.Sprintf("Operation %s:\n", alignText(opID.Name, 4)))
		var successful, failed, maxttc int
		for _, thread := range b.benchThreads {
			successful += thread.successful(opIDIndex)
			failed += thread.failed(opIDIndex)
			maxttc = int(math.Max(float64(maxttc), float64(thread.maxTTC(opIDIndex))))
		}

		printOutAndErr(fmt.Sprintf("  successful: %d\n  failed: %d\n  maxTTC: %d\n\n", successful, failed, maxttc))
		opType := opID.OpType
		opType.SuccessfulOperations += successful
		opType.FailedOperations += failed
		opType.Maxttc = int(math.Max(float64(maxttc), float64(opType.Maxttc)))
	}

	printOutAndErr(section("Sample errors (operation ratios [%])"))

	var totalSuccessful, totalFailed int
	for _, opType := range operations.OperationTypes {
		totalSuccessful += opType.SuccessfulOperations
		totalFailed += opType.FailedOperations
	}

	var totalError, totalTError float64
	for _, opType := range operations.OperationTypes {
		expectedRatio := opType.Probability * float64(opType.Count) * 100

		realRatio := float64(opType.SuccessfulOperations) / float64(totalSuccessful) * 100
		ratioError := math.Abs(expectedRatio - realRatio)

		tRealRatio := float64(opType.SuccessfulOperations+opType.FailedOperations) / float64(totalSuccessful+totalFailed) * 100
		tRatioError := math.Abs(expectedRatio - tRealRatio)

		totalError += ratioError
		totalTError += tRatioError

		printOutAndErr(fmt.Sprintf("%s:\n  expected: %s\n  successful: %s\n  error: %s\n  (total: %s\n  error: %s)\n", opType.Name, formatFloat(expectedRatio), formatFloat(realRatio), formatFloat(ratioError), formatFloat(tRealRatio), formatFloat(tRatioError)))
	}

	printOutAndErr(section("Summary results"))

	total := totalSuccessful + totalFailed
	for _, opType := range operations.OperationTypes {
		totalTypeOperations := opType.SuccessfulOperations + opType.FailedOperations
		printOutAndErr(fmt.Sprintf("%s:\n  successful: %d\n  maxTTC: %d\n  failed: %d\n  total: %d\n", opType.Name, opType.SuccessfulOperations, opType.Maxttc, opType.FailedOperations, totalTypeOperations))
	}
	printOutAndErr("\n")

	printOutAndErr(fmt.Sprintf("Total sample error: %s%% (%s%% including failed)\n", formatFloat(totalError), formatFloat(totalTError)))
	printOutAndErr(fmt.Sprintf("Total throughput: %s op/s (%s op/s including failed)\n", formatFloat(float64(totalSuccessful)/b.elapsedTime.Seconds()), formatFloat(float64(total)/b.elapsedTime.Seconds())))
	printOutAndErr(fmt.Sprintf("Elapsed time: %s\n", formatFloat(b.elapsedTime.Seconds())))
}

func logTTCIndexToTTC(logTTCIndex int) int {
	return int(float64(internal.MaxLowTTC+1) * math.Pow(internal.HighTTCLogBase, float64(logTTCIndex)))
}
