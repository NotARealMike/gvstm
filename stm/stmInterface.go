package stm

type TVar interface{}

type Transaction interface {
	Load(tVar TVar) interface{}
	Store(tVar TVar, value interface{})
}
