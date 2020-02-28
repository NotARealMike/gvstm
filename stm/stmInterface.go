package stm

type TVar interface {}

type Transaction interface {
	//commit() bool
	Load(tVar TVar) interface{}
	Store(tVar TVar, value interface{})
}

