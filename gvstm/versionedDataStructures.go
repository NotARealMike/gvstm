package gvstm

type vBox struct {
	body *vBody
}

func newVBox(initial interface{}) *vBox {
	vb := &vBox{}
	vb.body = &vBody{initial, 0, nil}
	return vb
}

func (vb *vBox) load(seqNo uint64) interface{} {
	body := vb.body
	for body.seqNo > seqNo {
		body = body.prev
	}
	return body.value
}

func (vb *vBox) commit(seqNo uint64, body *vBody) {
	prev := vb.body
	body.prev = prev
	body.seqNo = seqNo
	vb.body = body
}

type vBody struct {
	value interface{}
	seqNo uint64
	prev *vBody
}

