package interfaces

type ImmutableCollectionImpl struct {
	elements []interface{}
}

func NewImmutableCollectionImpl(elements []interface{}) *ImmutableCollectionImpl {
	ic := &ImmutableCollectionImpl{make([]interface{}, len(elements))}
	for i := range elements {
		ic.elements[i] = elements[i]
	}
	return ic
}

func (ic *ImmutableCollectionImpl) Size() int {
	return len(ic.elements)
}

func (ic *ImmutableCollectionImpl) Contains(element interface{}) bool {
	for _, e := range ic.elements {
		if e == element {
			return true
		}
	}
	return false
}

func (ic *ImmutableCollectionImpl) Clone() ImmutableCollection {
	return NewImmutableCollectionImpl(ic.elements)
}

func (ic *ImmutableCollectionImpl) ToSlice() []interface{} {
	return ic.elements
}
