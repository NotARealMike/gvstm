package interfaces

import . "gvstm/stm"

// Backend data structures used by the benchmark data structures.
// An implementation of the benchmark must provide correctly synchronised
// implementations of these interfaces.

// A pool of ids to generate a unique id for each instance of each
// type in the benchmark data structures.
type IDPool interface {
	GetID(tx Transaction) (int, error)
	PutUnusedID(tx Transaction, id int)
}

type ImmutableCollection interface {
	Size() int
	Contains(element interface{}) bool
	Clone() ImmutableCollection
	ToSlice() []interface{}
}

// Indexes are used by many of the benchmark operations.
type Index interface {
	Get(tx Transaction, key interface{}) interface{}
	Put(tx Transaction, key interface{}, value interface{})
	PutIfAbsent(tx Transaction, key interface{}, value interface{})
	Remove(tx Transaction, key interface{}) bool
	GetRange(tx Transaction, minKey interface{}, maxKey interface{}) []interface{}
	GetKeys(tx Transaction) []interface{}
}

// A set with at least a few hundred elements.
// Used to hold the atomic parts of a composite part and
// the index of atomic parts with the same build date.
type LargeSet interface {
	Add(tx Transaction, element interface{}) bool
	Remove(tx Transaction, element interface{}) bool
	Contains(tx Transaction, element interface{}) bool
	Size(tx Transaction) int
	ToSlice(tx Transaction) []interface{}
}
