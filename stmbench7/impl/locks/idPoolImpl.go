package locks

import (
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
)

type idPoolImpl struct {
	pool []int
	head int
}

func newIDPoolImpl(tx Transaction, maxNumberOfIDs int) IDPool {
	pool := make([]int, maxNumberOfIDs)
	for i := range pool {
		pool[i] = i + 1
	}
	return &idPoolImpl{
		pool: pool,
		head: 0,
	}
}

func (ip *idPoolImpl) GetID(tx Transaction) (int, OpFailedError) {
	if ip.head == len(ip.pool) {
		return -1, NewOpFailedError("IDPool exhausted")
	}
	ip.head++
	return ip.pool[ip.head-1], nil
}

func (ip *idPoolImpl) PutUnusedID(tx Transaction, id int) {
	oldPool := ip.pool
	newPool := make([]int, len(oldPool))
	for i := range oldPool {
		newPool[i] = oldPool[i]
	}
	newHead := ip.head - 1
	newPool[newHead] = id
	ip.pool = newPool
	ip.head = newHead
}
