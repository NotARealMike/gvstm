package gvstm

import "gvstm/stm"

func TVarFactory(value interface{}) stm.TVar {
	vb := newVBox(value)
	return vb
}