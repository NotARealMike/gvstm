package gvstm

import (
    "gvstm/gvstm"
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type indexImpl struct {
    index TVar
}

func newIndexImpl(tx Transaction) Index {
    return &indexImpl{
        index:gvstm.CreateTVar(map[interface{}]interface{}{}),
    }
}

func (i *indexImpl) Get(tx Transaction, key interface{}) interface{} {
    index := tx.Load(i.index).(map[interface{}]interface{})
    return index[key]
}

func (i *indexImpl) Put(tx Transaction, key interface{}, value interface{}) {
    // TODO: Put-ing a nil value should panic
    oldIndex := tx.Load(i.index).(map[interface{}]interface{})
    newIndex := make(map[interface{}]interface{}, len(oldIndex))
    for k, v := range oldIndex {
        newIndex[k] = v
    }
    newIndex[key] = value
    tx.Store(i.index, newIndex)
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
    oldIndex := tx.Load(i.index).(map[interface{}]interface{})
    newIndex := make(map[interface{}]interface{}, len(oldIndex))
    for k, v := range oldIndex {
        newIndex[k] = v
    }
    delete(newIndex, key)
    tx.Store(i.index, newIndex)
    return true
}

func (i *indexImpl) GetRange(tx Transaction, minKey interface{}, maxKey interface{}) []interface{} {
    min := minKey.(int)
    max := maxKey.(int)
    index := tx.Load(i.index).(map[interface{}]interface{})
    values := make([]interface{}, 0, len(index))
    for k, v := range index {
        if ki := k.(int); ki >= min && ki <= max {
            values = append(values, v)
        }
    }
    return values
}

func (i *indexImpl) GetKeys(tx Transaction) []interface{} {
    index := tx.Load(i.index).(map[interface{}]interface{})
    keys := make([]interface{}, len(index))
    j := 0
    for k := range index {
        keys[j] = k
        j++
    }
    return keys
}
