package gvstm

type VBox struct {
	body *vBody
}

func NewVBox(initial interface{}) *VBox {
	vb := &VBox{}
	vb.body = &vBody{initial, 0, nil}
	return vb
}

func (vb *VBox) read(seqNo uint64) interface{} {
	body := vb.body
	for body.seqNo > seqNo {
		body = body.prev
	}
	return body.value
}

func (vb *VBox) commit(seqNo uint64, body *vBody) {
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

