package gvstm

import (
	"github.com/NotARealMike/gvstm/stm"
	"sync"
	"unsafe"
)

var (
	commitMutex sync.Mutex
)

type rWTransaction struct {
	latestRecord *activeTxRecord
	readSet      map[*vBox]struct{}
	writeSet     map[*vBox]*vBody
}

func newRWTransaction() *rWTransaction {
	return &rWTransaction{
		latestTxRec.registerTransaction(),
		map[*vBox]struct{}{},
		map[*vBox]*vBody{},
	}
}

func (tx *rWTransaction) commit() bool {
	if len(tx.writeSet) == 0 {
		return true
	}

	commitMutex.Lock()
	// Validation
	for vb := range tx.readSet {
		if vb.body.Load().(vBody).seqNo > tx.latestRecord.txNumber {
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

func (tx *rWTransaction) Load(tVar stm.TVar) interface{} {
	vb := tVar.(*vBox)
	if body, ok := tx.writeSet[vb]; ok {
		return body.value
	}
	tx.readSet[vb] = struct{}{}
	return vb.load(tx.latestRecord.txNumber)
}

func (tx *rWTransaction) Store(tVar stm.TVar, value interface{}) {
	// The sequence number and prev pointer will be set to the correct values at commit time.
	// The vBody is created here to avoid memory allocation when holding the commit lock.
	vb := tVar.(*vBox)
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

func (tx *readOnlyTransaction) Load(tVar stm.TVar) interface{} {
	vb := tVar.(*vBox)
	return vb.load(tx.latestRecord.txNumber)
}

func (tx *readOnlyTransaction) Store(tVar stm.TVar, value interface{}) {
	panic("write in a read only transaction")
}

func Atomic(f func(tx stm.Transaction)) {
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
	// Load only transactions always commit.
	readTx = newReadOnlyTransaction()
	f(readTx)
	readTx.commit()
	return
}
