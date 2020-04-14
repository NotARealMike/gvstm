package gvstm

import "gvstm/stm"

func CreateTVar(value interface{}) stm.TVar {
	vb := newVBox(value)
	return vb
}