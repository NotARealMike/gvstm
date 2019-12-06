package gvstm

import (
	"sync"
	"sync/atomic"
)

var (
	globalClock clock
	commitMutex sync.Mutex
)

type clock uint64

func (c *clock) updateClock(t uint64) {
	atomic.StoreUint64((*uint64)(c), t)
}

func (c *clock) sampleClock() (t uint64) {
	return atomic.LoadUint64((*uint64)(c))
}

type VBox struct {
	body atomic.Value
}

func NewVBox(initial interface{}) *VBox {
	body := vBody{initial, 0, nil}
	vb := &VBox{}
	vb.body.Store(body)
	return vb
}

func (vb *VBox) read(seqNo uint64) interface{} {
	body := vb.body.Load().(vBody)
	for body.seqNo > seqNo {
		body = *body.prev
	}
	return body.value
}

func (vb *VBox) commit(seqNo uint64, value interface{}) {
	prev := vb.body.Load().(vBody)
	body := vBody{value, seqNo, &prev}
	vb.body.Store(body)
}

type vBody struct {
	value interface{}
	seqNo uint64
	prev *vBody
}

type Transaction interface {
	commit() bool
	Read(vb *VBox) interface{}
	Write(vb *VBox, value interface{})
}

type rWTransaction struct {
	readVersion uint64
	readSet map[*VBox]struct{}
	writeSet map[*VBox]interface{}
}

func newRWTransaction() *rWTransaction {
	return &rWTransaction{
		globalClock.sampleClock(),
		map[*VBox]struct{}{},
		map[*VBox]interface{}{},
	}
}

func (tx *rWTransaction) commit() bool {
	if len(tx.writeSet) == 0 {
		return true
	}
	commitMutex.Lock()
	for vb := range tx.readSet {
		if vb.body.Load().(vBody).seqNo > tx.readVersion {
			commitMutex.Unlock()
			return false
		}
	}
	writeVersion := globalClock.sampleClock() + 1
	for vb, value := range tx.writeSet {
		vb.commit(writeVersion, value)
	}
	globalClock.updateClock(writeVersion)
	commitMutex.Unlock()
	return true
}

func (tx *rWTransaction) Read(vb *VBox) interface{} {
	if value, ok := tx.writeSet[vb]; ok {
		return value
	}
	tx.readSet[vb] = struct {}{}
	return vb.read(tx.readVersion)
}

func (tx *rWTransaction) Write(vb *VBox, value interface{}) {
	tx.writeSet[vb] = value
}

type readOnlyTransaction struct {
	readVersion uint64
}

func newReadOnlyTransaction() *readOnlyTransaction {
	return &readOnlyTransaction{globalClock.sampleClock()}
}

func (tx *readOnlyTransaction) commit() bool {
	return true
}

func (tx *readOnlyTransaction) Read(vb *VBox) interface{} {
	return vb.read(tx.readVersion)
}

func (tx *readOnlyTransaction) Write(vb *VBox, value interface{}) {
	panic("write in a read only transaction")
}

func Atomic(f func(tx Transaction)) {
	var tx Transaction

	defer func() {
		// When a read only transaction tries to write it panics and ends up here.
		if r := recover(); r != nil {
		RW:
			tx = newRWTransaction()
			f(tx)
			if !tx.commit() {
				goto RW
			}
			return
		}
	}()

	// All transactions start out as read only transactions.
	// Read only transactions always commit.
	tx = newReadOnlyTransaction()
	f(tx)
	return
}
