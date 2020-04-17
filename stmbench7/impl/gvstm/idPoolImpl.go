package gvstm

import (
    "gvstm/gvstm"
    . "gvstm/stm"
    . "gvstm/stmbench7/interfaces"
)

type idPoolImpl struct {
    pool TVar
    head TVar
}

func newIDPoolImpl(tx Transaction, maxNumberOfIDs int) IDPool {
    pool := make([]int, maxNumberOfIDs)
    for i := range pool {
        pool[i] = i+1
    }
    return &idPoolImpl{
        pool: gvstm.CreateTVar(pool),
        head: gvstm.CreateTVar(0),
    }
}

func (ip *idPoolImpl) GetID(tx Transaction) (int, OpFailedError) {
    pool := tx.Load(ip.pool).([]int)
    head := tx.Load(ip.head).(int)
    if head == len(pool) {
        return -1, NewOpFailedError("IDPool exhausted")
    }
    tx.Store(ip.head, head+1)
    return pool[head], nil
}

func (ip *idPoolImpl) PutUnusedID(tx Transaction, id int) {
    oldPool := tx.Load(ip.pool).([]int)
    newPool := make([]int, len(oldPool))
    for i := range oldPool {
        newPool[i] = oldPool[i]
    }
    newHead := tx.Load(ip.head).(int) - 1
    newPool[newHead] = id
    tx.Store(ip.pool, newPool)
    tx.Store(ip.head, newHead)
}
