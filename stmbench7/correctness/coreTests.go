package correctness

import (
	"fmt"
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"math"
	"os"
	"strings"
)

func designObjInvariantTest(tx Transaction, obj DesignObj, initial bool, maxID, minBuildDate, maxBuildDate int) error {
	id := obj.GetId(tx)
	if id < 0 || id > maxID {
		return newOutOfBoundError(obj, id, "id", 1, maxID, id)
	}
	if !checkValidType(obj.GetType(tx)) {
		return newUnexpectedStringError(obj, id, "type", "type #(num)", obj.GetType(tx))
	}
	buildDate := obj.GetBuildDate(tx)
	if !initial {
		if minBuildDate%2 == 0 {
			minBuildDate--
		}
		if maxBuildDate%2 != 0 {
			maxBuildDate++
		}
	}
	if buildDate < minBuildDate || buildDate > maxBuildDate {
		return newOutOfBoundError(obj, id, "buildDate", minBuildDate, maxBuildDate, buildDate)
	}
	return nil
}

func moduleInvariantTest(tx Transaction, module Module, initial bool, objects traversedObjects) error {
	if err := designObjInvariantTest(tx, module, initial, internal.NumModules, internal.MinModuleDate, internal.MaxModuleDate); err != nil {
		return err
	}
	id := module.GetId(tx)
	manual := module.GetManual(tx)
	if manual == nil {
		return newInvFailedError(module, id, "Null manual in module.")
	}
	if err := manualInvariantTest(tx, manual, module); err != nil {
		return newInvFailedError(module, id, err.Error())
	}
	rootAssembly := module.GetDesignRoot(tx)
	if rootAssembly == nil {
		return newInvFailedError(module, id, "Null root assembly.")
	}
	if err := complexAssemblyInvariantTest(tx, rootAssembly, nil, initial, module, objects); err != nil {
		return newInvFailedError(module, id, err.Error())
	}
	return nil
}

func manualInvariantTest(tx Transaction, manual Manual, module Module) error {
	id := manual.GetId(tx)
	if id != 1 {
		return newUnexpectedIntError(manual, id, "id", 1, id)
	}
	title := manual.GetTitle(tx)
	titleShouldBe := "Manual for module #1"
	if title != titleShouldBe {
		return newUnexpectedStringError(manual, id, "title", titleShouldBe, title)
	}
	if !(manual.StartsWith(tx, 'I') || manual.StartsWith(tx, 'i')) {
		return newUnexpectedStringError(manual, id, "text (prefix)", "'I' || 'i'", manual.GetText(tx)[0:7]+"...")
	}
	if manual.StartsWith(tx, 'I') && manual.CountOccurences(tx, 'i') > 0 {
		return newInvFailedError(manual, id, "text starts from 'I' but contains 'i'")
	}
	if manual.StartsWith(tx, 'i') && manual.CountOccurences(tx, 'I') > 0 {
		return newInvFailedError(manual, id, "text starts from 'i' but contains 'I'")
	}
	if manual.GetModule(tx) != module {
		return newInvFailedError(manual, id, "invalid connection to module")
	}
	return nil
}

func assemblyInvariantTest(tx Transaction, assembly Assembly, parent ComplexAssembly, initial bool, maxID int, module Module) error {
	if err := designObjInvariantTest(tx, assembly, initial, maxID, internal.MinAssDate, internal.MaxAssDate); err != nil {
		return err
	}
	id := assembly.GetId(tx)
	if assembly.GetSuperAssembly(tx) != parent {
		return newInvFailedError(assembly, id, "invalid reference to the parent ComplexAssembly")
	}
	if assembly.GetModule(tx) != module {
		return newInvFailedError(assembly, id, "invalid reference to the parent Module")
	}
	return nil
}

func complexAssemblyInvariantTest(tx Transaction, assembly, parentAssembly ComplexAssembly, initial bool, module Module, objects traversedObjects) error {
	objects.complexAssemblies[assembly] = struct{}{}
	if err := assemblyInvariantTest(tx, assembly, parentAssembly, initial, internal.MaxComplexAssemblies, module); err != nil {
		return err
	}

	id := assembly.GetId(tx)
	level := assembly.GetLevel(tx)
	if level <= 1 || level > internal.NumAssemblyLevels {
		return newOutOfBoundError(assembly, id, "level", 1, internal.NumAssemblyLevels, level)
	}
	for _, subAssembly := range assembly.GetSubAssemblies(tx).ToSlice() {
		if level > 2 {
			if _, ok := subAssembly.(ComplexAssembly); !ok {
				return newInvFailedError(assembly, id, fmt.Sprintf("subAssembly not of type ComplexAssembly at level = %d", level))
			}
			if err := complexAssemblyInvariantTest(tx, subAssembly.(ComplexAssembly), assembly, initial, module, objects); err != nil {
				return newInvFailedError(assembly, id, err.Error())
			}
		} else {
			if _, ok := subAssembly.(BaseAssembly); !ok {
				return newInvFailedError(assembly, id, "subAssembly not of type BaseAssembly at level = 2")
			}
			if err := baseAssemblyInvariantTest(tx, subAssembly.(BaseAssembly), assembly, initial, module, objects); err != nil {
				return newInvFailedError(assembly, id, err.Error())
			}
		}
	}
	return nil
}

func baseAssemblyInvariantTest(tx Transaction, assembly BaseAssembly, parent ComplexAssembly, initial bool, module Module, objects traversedObjects) error {
	objects.baseAssemblies[assembly] = struct{}{}
	if err := assemblyInvariantTest(tx, assembly, parent, initial, internal.MaxBaseAssemblies, module); err != nil {
		return err
	}
	for _, compositePart := range assembly.GetComponents(tx).ToSlice() {
		if err := compositePartInvariantTest(tx, compositePart.(CompositePart), initial, assembly, objects); err != nil {
			return newInvFailedError(assembly, assembly.GetId(tx), err.Error())
		}
	}
	return nil
}

func compositePartInvariantTest(tx Transaction, component CompositePart, initial bool, parent BaseAssembly, objects traversedObjects) error {
	objects.components[component] = struct{}{}
	fmt.Fprintf(os.Stderr, "Component %d\n", len(objects.components))
	if err := designObjInvariantTest(tx, component, initial, internal.MaxCompParts, 0, 1073741824); err != nil {
		return err
	}

	id := component.GetId(tx)
	buildDate := component.GetBuildDate(tx)
	minOldDate, maxOldDate, minYoungDate, maxYoungDate := internal.MinOldCompositePartDate, internal.MaxOldCompositePartDate, internal.MinYoungCompositePartDate, internal.MaxYoungCompositePartDate
	if !initial {
		if minOldDate%2 == 0 {
			minOldDate--
		}
		if maxOldDate%2 != 0 {
			maxOldDate++
		}
		if minYoungDate%2 == 0 {
			minYoungDate--
		}
		if maxYoungDate%2 != 0 {
			maxYoungDate++
		}
	}
	if (buildDate < minOldDate || buildDate > maxOldDate) && (buildDate < minYoungDate || buildDate > maxYoungDate) {
		return newInvFailedError(component, id, fmt.Sprintf("wrong buildDate (%d)", buildDate))
	}

	usedIn := component.GetUsedIn(tx)
	if parent != nil && !usedIn.Contains(parent) {
		return newInvFailedError(component, id, "a BaseAssembly parent not in set usedIn")
	}

	for _, assembly := range usedIn.ToSlice() {
		if !assembly.(BaseAssembly).GetComponents(tx).Contains(component) {
			return newInvFailedError(component, id, newInvFailedError(assembly, assembly.(BaseAssembly).GetId(tx), "a child CompositePart not in set of components").Error())
		}
	}

	if err := documentInvariantTest(tx, component.GetDocumentation(tx), component, objects); err != nil {
		return newInvFailedError(component, id, err.Error())
	}

	parts := component.GetParts(tx)
	if !parts.Contains(tx, component.GetRootPart(tx)) {
		return newInvFailedError(component, id, "rootPart not in set of parts")
	}

	allConnectedParts := map[AtomicPart]struct{}{}
	for _, part := range parts.ToSlice(tx) {
		if err := atomicPartInvariantTest(tx, part.(AtomicPart), component, initial, allConnectedParts, objects); err != nil {
			return newInvFailedError(component, id, err.Error())
		}
	}
	for part := range allConnectedParts {
		if !parts.Contains(tx, part) {
			return newInvFailedError(component, id, "a graph-connected AtomicPart not in set of all parts")
		}
	}
	return nil
}

func atomicPartInvariantTest(tx Transaction, part AtomicPart, component CompositePart, initial bool, allParts map[AtomicPart]struct{}, objects traversedObjects) error {
	objects.atomicParts[part] = struct{}{}
	if err := designObjInvariantTest(tx, part, initial, internal.MaxAtomicParts, internal.MinAtomicPartDate, internal.MaxAtomicPartDate); err != nil {
		return err
	}

	id := part.GetId(tx)
	x := part.GetX(tx)
	y := part.GetY(tx)
	if math.Abs(float64(x-y)) != 1 {
		return newInvFailedError(part, id, fmt.Sprintf("inconsistent x and y attributes: x = %d, y = %d", x, y))
	}

	if part.GetPartOf(tx) != component {
		return newInvFailedError(part, id, "invalid reference to CompositePart parent")
	}

	to := part.GetToConnections(tx)
	from := part.GetFromConnections(tx)
	if to.Size() != internal.NumConnectionsPerAtomicPart {
		return newUnexpectedIntError(part, id, "to", internal.NumConnectionsPerAtomicPart, to.Size())
	}

	for _, connection := range from.ToSlice() {
		if err := connectionInvariantTest(connection.(Connection), part); err != nil {
			return newInvFailedError(part, id, err.Error())
		}
		sourceToSet := connection.(Connection).GetSource().GetToConnections(tx)
		foundMyself := false
		for _, outConn := range sourceToSet.ToSlice() {
			if outConn.(Connection).GetSource() == part {
				foundMyself = true
				break
			}
		}
		if !foundMyself {
			return newInvFailedError(part, id, "inconsistent 'from' set")
		}
	}

	for _, connection := range to.ToSlice() {
		if err := connectionInvariantTest(connection.(Connection), part); err != nil {
			return newInvFailedError(part, id, err.Error())
		}
		allParts[connection.(Connection).GetDestination()] = struct{}{}
	}

	return nil
}

func connectionInvariantTest(connection Connection, from AtomicPart) error {
	if !checkValidType(connection.GetType()) {
		return newUnexpectedStringError(connection, 0, "type", "type #...", connection.GetType())
	}
	length := connection.GetLength()
	if length < 1 || length > internal.XYRange {
		return newOutOfBoundError(connection, 0, "length", 1, internal.XYRange, length)
	}
	if connection.GetSource() != from {
		return newInvFailedError(connection, 0, "invalid source (from) reference")
	}
	return nil
}

func documentInvariantTest(tx Transaction, document Document, component CompositePart, objects traversedObjects) error {
	objects.documents[document] = struct{}{}

	id := document.GetDocumentId(tx)
	if id < 1 || id > internal.MaxCompParts {
		return newOutOfBoundError(document, id, "id", 1, internal.MaxCompParts, id)
	}
	if document.GetCompositePart(tx) != component {
		return newInvFailedError(document, id, "invalid reference to CompositePart parent")
	}

	title := document.GetTitle(tx)
	titleShouldBe := fmt.Sprintf("Composite Part #%d", component.GetId(tx))
	if title != titleShouldBe {
		return newUnexpectedStringError(document, id, "title", titleShouldBe, title)
	}

	text := document.GetText(tx)
	if !(strings.HasPrefix(text, fmt.Sprintf("I am the documentation for composite part #%d", component.GetId(tx))) ||
		strings.HasPrefix(text, fmt.Sprintf("This is the documentation for composite part #%d", component.GetId(tx)))) {
		return newUnexpectedStringError(document, id, "text (prefx)", "I am / This is the documentation for composite part #...", text[0:30])
	}
	return nil
}
