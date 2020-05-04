package locks

import (
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
)

type indexImpl struct {
	index map[interface{}]interface{}
}

func newIndexImpl(tx Transaction) Index {
	return &indexImpl{
		index: map[interface{}]interface{}{},
	}
}

func (i *indexImpl) Get(tx Transaction, key interface{}) interface{} {
	return i.index[key]
}

func (i *indexImpl) Put(tx Transaction, key interface{}, value interface{}) {
	newIndex := make(map[interface{}]interface{}, len(i.index))
	for k, v := range i.index {
		newIndex[k] = v
	}
	newIndex[key] = value
	i.index = newIndex
}

func (i *indexImpl) PutIfAbsent(tx Transaction, key interface{}, value interface{}) interface{} {
	if oldValue := i.Get(tx, key); oldValue != nil {
		return oldValue
	}
	i.Put(tx, key, value)
	return nil
}

func (i *indexImpl) Remove(tx Transaction, key interface{}) bool {
	if i.Get(tx, key) == nil {
		return false
	}
	newIndex := make(map[interface{}]interface{}, len(i.index))
	for k, v := range i.index {
		newIndex[k] = v
	}
	delete(newIndex, key)
	i.index = newIndex
	return true
}

func (i *indexImpl) GetRange(tx Transaction, minKey interface{}, maxKey interface{}) []interface{} {
	min := minKey.(int)
	max := maxKey.(int)
	values := make([]interface{}, 0, len(i.index))
	for k, v := range i.index {
		if ki := k.(int); ki >= min && ki <= max {
			values = append(values, v)
		}
	}
	return values
}

func (i *indexImpl) GetKeys(tx Transaction) []interface{} {
	keys := make([]interface{}, len(i.index))
	j := 0
	for k := range i.index {
		keys[j] = k
		j++
	}
	return keys
}
