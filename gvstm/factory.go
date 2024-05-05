package gvstm

import "github.com/NotARealMike/gvstm/stm"

func CreateTVar(value interface{}) stm.TVar {
	vb := newVBox(value)
	return vb
}
