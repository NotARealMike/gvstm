package interfaces

import . "gvstm/stm"

// Interfaces of the benchmark data structures.

type Assembly interface {
	DesignObj
	GetSuperAssembly(tx Transaction) ComplexAssembly
	GetModule(tx Transaction) Module
	ClearPointers(tx Transaction)
}

type AtomicPart interface {
	DesignObj
	ConnectTo(tx Transaction, destination AtomicPart, typ string, length int)
	AddConnectionFromOtherPart(tx Transaction, connection Connection)
	SetCompositePart(tx Transaction, partOf CompositePart)
	GetNumToConnections(tx Transaction) int
	GetToConnections(tx Transaction) ImmutableCollection
	GetFromConnections(tx Transaction) ImmutableCollection
	GetPartOf(tx Transaction) CompositePart
	SwapXY(tx Transaction)
	GetX(tx Transaction) int
	GetY(tx Transaction) int
	ClearPointers(tx Transaction)
}

type BaseAssembly interface {
	Assembly
	AddComponent(tx Transaction, component CompositePart)
	RemoveComponent(tx Transaction, component CompositePart) bool
	GetComponents(tx Transaction) ImmutableCollection
}

type ComplexAssembly interface {
	Assembly
	AddSubAssembly(tx Transaction, assembly Assembly) bool
	RemoveSubAssembly(tx Transaction, assembly Assembly) bool
	GetSubAssemblies(tx Transaction) ImmutableCollection
	GetLevel(tx Transaction) int
}

type CompositePart interface {
	DesignObj
	AddAssembly(tx Transaction, assembly BaseAssembly)
	AddPart(tx Transaction, part AtomicPart) bool
	SetRootPart(tx Transaction, part AtomicPart)
	GetRootPart(tx Transaction) AtomicPart
	GetDocumentation(tx Transaction) Document
	GetParts(tx Transaction) LargeSet
	RemoveAssembly(tx Transaction, assembly BaseAssembly)
	GetUsedIn(tx Transaction) ImmutableCollection
	ClearPointers(tx Transaction)
}

// Connections are immutable, so a single,
// non-synchronised implementation is sufficient.
type Connection interface {
	GetReversed() Connection
	GetDestination() AtomicPart
	GetSource() AtomicPart
	GetType() string
	GetLength() int
}

type DesignObj interface {
	GetId(tx Transaction) int
	GetBuildDate(tx Transaction) int
	GetType(tx Transaction) string
	UpdateBuildDate(tx Transaction)
	NullOperation(tx Transaction)
}

type Document interface {
	SetPart(tx Transaction, part CompositePart)
	GetCompositePart(tx Transaction) CompositePart
	GetDocumentId(tx Transaction) int
	GetTitle(tx Transaction) string
	NullOperation(tx Transaction)
	SearchText(tx Transaction, symbol rune) int
	ReplaceText(tx Transaction, from string, to string) int
	TextBeginsWith(tx Transaction, prefix string) bool
	GetText(tx Transaction) string
}

type Manual interface {
	GetId(tx Transaction) int
	GetTitle(tx Transaction) string
	GetText(tx Transaction) string
	GetModule(tx Transaction) Module
	SetModule(tx Transaction, module Module)
	CountOccurences(tx Transaction, ch rune) int
	CheckFirstLastCharTheSame(tx Transaction) int
	StartsWith(tx Transaction, ch rune) bool
	ReplaceChar(tx Transaction, from rune, to rune) int
}

type Module interface {
	DesignObj
	SetDesignRoot(tx Transaction, designRoot ComplexAssembly)
	GetDesignRoot(tx Transaction) ComplexAssembly
	GetManual(tx Transaction) Manual
}
