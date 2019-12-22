package gvstm

import "sync/atomic"

type VBox struct {
	// Stores a pointer to a vBody.
	body atomic.Value
}

func NewVBox(initial interface{}) *VBox {
	body := &vBody{initial, 0, nil}
	vb := &VBox{}
	vb.body.Store(body)
	return vb
}

func (vb *VBox) read(seqNo uint64) interface{} {
	body := vb.body.Load().(*vBody)
	for body.seqNo > seqNo {
		body = body.prev
	}
	return body.value
}

func (vb *VBox) commit(seqNo uint64, body *vBody) {
	prev := vb.body.Load().(*vBody)
	body.prev = prev
	body.seqNo = seqNo
	vb.body.Store(body)
}

type vBody struct {
	value interface{}
	seqNo uint64
	prev *vBody
}

