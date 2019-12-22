package gvstm

import (
	"sync/atomic"
	"unsafe"
)

var (
	latestTxRec *activeTxRecord = &activeTxRecord{
		0,
		unsafe.Pointer(nil),
		0,
		nil,
	}
)

type activeTxRecord struct {
	// The commit timestamp.
	txNumber uint64
	// Holds the bodies that need to be cleaned when all older transactions finish.
	// Needs to be an unsafe pointer so that a CAS can be performed.
	// The type of the value stored is []*vBody.
	bodiesToClean unsafe.Pointer
	// The number of active transactions reading with this timestamp.
	running uint64
	// The next active transaction record in the list.
	next *activeTxRecord
}

func (atr *activeTxRecord) incrementRunning() {
	atomic.AddUint64(&(atr.running), 1)
}

func (atr *activeTxRecord) decrementRunning() {
	// The atomic package's recommended way of decrementing.
	if atomic.AddUint64(&atr.running, ^uint64(0)) == 0 {
		atr.maybeCleanSuc()
	}
}

func (atr *activeTxRecord) maybeCleanSuc() {
	for true {
		if atr.next != nil && atr.bodiesToClean == nil && atr.running == 0 && atr.next.clean() {
			atr = atr.next
		} else {
			break
		}
	}
}

func (atr *activeTxRecord) clean() bool {
	unsafeBodies := atomic.SwapPointer(&atr.bodiesToClean, nil)

	if unsafeBodies == nil {
		return false
	}

	bodies := *(*[]*vBody)(unsafeBodies)
	for _, body := range bodies {
		body.prev = nil
	}
	return true
}

func (atr *activeTxRecord) registerTransaction() *activeTxRecord {
	for true {
		atr.incrementRunning()
		if atr.next == nil {
			return atr
		} else {
			atr.decrementRunning()
			atr = atr.next
		}
	}
	return nil
}
