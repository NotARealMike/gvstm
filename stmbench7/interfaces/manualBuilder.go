package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
    "strconv"
)

type ManualBuilder struct {
    idPool IDPool
}

func NewManualBuilder(tx Transaction) *ManualBuilder {
    return &ManualBuilder{
        idPool: beFactory.CreateIDPool(tx, internal.NumModules),
    }
}

func (mb *ManualBuilder) CreateManual(tx Transaction, moduleID int) (Manual, *OpFailedError) {
    manualID, err := mb.idPool.GetID(tx)
    if err != nil {
        return nil, err
    }
    title := "Manual for module #" + strconv.Itoa(moduleID)
    text := createText(internal.ManualSize, "I am the manual for module #" + strconv.Itoa(moduleID) + "\n")
    return doFactory.CreateManual(tx, manualID, title, text), nil
}
