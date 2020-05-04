package interfaces

import (
	. "gvstm/stm"
	"gvstm/stmbench7/internal"
	"strconv"
)

type manualBuilder interface {
	createManual(tx Transaction, moduleID int) (Manual, OpFailedError)
}

type manualBuilderImpl struct {
	idPool IDPool
}

func newManualBuilder(tx Transaction) manualBuilder {
	return &manualBuilderImpl{
		idPool: BEFactory.CreateIDPool(tx, internal.NumModules),
	}
}

func (mb *manualBuilderImpl) createManual(tx Transaction, moduleID int) (Manual, OpFailedError) {
	manualID, err := mb.idPool.GetID(tx)
	if err != nil {
		return nil, err
	}
	title := "Manual for module #" + strconv.Itoa(moduleID)
	text := createText(internal.ManualSize, "I am the manual for module #"+strconv.Itoa(moduleID)+"\n")
	return DOFactory.CreateManual(tx, manualID, title, text), nil
}
