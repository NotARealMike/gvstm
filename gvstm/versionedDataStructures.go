package gvstm

import "sync/atomic"

type vBox struct {
	body atomic.Value
}

func newVBox(initial interface{}) *vBox {
	vb := &vBox{}
	vb.body.Store(vBody{initial, 0, nil})
	return vb
}

func (vb *vBox) load(seqNo uint64) interface{} {
	body := vb.body.Load().(vBody)
	bp := &body
	for bp.seqNo > seqNo {
		bp = bp.prev
	}
	return bp.value
}

func (vb *vBox) commit(seqNo uint64, body *vBody) {
	prev := vb.body.Load().(vBody)
	body.prev = &prev
	body.seqNo = seqNo
	vb.body.Store(*body)
}

type vBody struct {
	value interface{}
	seqNo uint64
	prev  *vBody
}
