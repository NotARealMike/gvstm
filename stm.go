package gvstm

import (
	"sync"
	"unsafe"
)

var (
	commitMutex sync.Mutex
)

type Transaction interface {
	commit() bool
	Read(vb *VBox) interface{}
	Write(vb *VBox, value interface{})
}

type rWTransaction struct {
	latestRecord *activeTxRecord
	readSet map[*VBox]struct{}
	writeSet map[*VBox]*vBody
}

func newRWTransaction() *rWTransaction {
	return &rWTransaction{
		latestTxRec.registerTransaction(),
		map[*VBox]struct{}{},
		map[*VBox]*vBody{},
	}
}

func (tx *rWTransaction) commit() bool {
	if len(tx.writeSet) == 0 {
		return true
	}

	commitMutex.Lock()
	// Validation
	for vb := range tx.readSet {
		if vb.body.seqNo > tx.latestRecord.txNumber {
			commitMutex.Unlock()
			return false
		}
	}

	// Commit
	writeVersion := latestTxRec.txNumber + 1
	bodies := make([]*vBody, len(tx.writeSet))
	index := 0
	for vb, body := range tx.writeSet {
		vb.commit(writeVersion, body)
		bodies[index] = body
		index++
	}
	newTxRecord := &activeTxRecord{
		txNumber:      writeVersion,
		bodiesToClean: unsafe.Pointer(&bodies),
		running:       1,
		next:          nil,
	}
	latestTxRec.next = newTxRecord
	oldTxRecord := latestTxRec
	latestTxRec = newTxRecord
	commitMutex.Unlock()

	oldTxRecord.decrementRunning()
	newTxRecord.decrementRunning()

	return true
}

func (tx *rWTransaction) Read(vb *VBox) interface{} {
	if body, ok := tx.writeSet[vb]; ok {
		return body.value
	}
	tx.readSet[vb] = struct{}{}
	return vb.read(tx.latestRecord.txNumber)
}

func (tx *rWTransaction) Write(vb *VBox, value interface{}) {
	// The sequence number and prev pointer will be set to the correct values at commit time.
	// The vBody is created here to avoid memory allocation when holding the commit lock.
	tx.writeSet[vb] = &vBody{
		value,
		0,
		nil,
	}
}

type readOnlyTransaction struct {
	latestRecord *activeTxRecord
}

func newReadOnlyTransaction() *readOnlyTransaction {
	return &readOnlyTransaction{latestTxRec.registerTransaction()}
}

func (tx *readOnlyTransaction) commit() bool {
	tx.latestRecord.decrementRunning()
	return true
}

func (tx *readOnlyTransaction) Read(vb *VBox) interface{} {
	return vb.read(tx.latestRecord.txNumber)
}

func (tx *readOnlyTransaction) Write(vb *VBox, value interface{}) {
	panic("write in a read only transaction")
}

func Atomic(f func(tx Transaction)) {
	var readTx *readOnlyTransaction

	defer func() {
		// When a read only transaction tries to write it panics and ends up here.
		if r := recover(); r != nil {
			readTx.latestRecord.decrementRunning()
		RW:
			rwTx := newRWTransaction()
			f(rwTx)
			if !rwTx.commit() {
				rwTx.latestRecord.decrementRunning()
				goto RW
			}
			return
		}
	}()

	// All transactions start out as read only transactions.
	// Read only transactions always commit.
	readTx = newReadOnlyTransaction()
	f(readTx)
	readTx.commit()
	return
}
