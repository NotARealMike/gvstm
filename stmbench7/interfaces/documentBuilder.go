package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
    "strconv"
)

type DocumentBuilder struct {
    idPool IDPool
    documentTitleIndex Index
}

func NewDocumentBuilder(tx Transaction, documentTitleIndex Index) *DocumentBuilder {
    return &DocumentBuilder{
        idPool:             beFactory.CreateIDPool(tx, internal.MaxCompParts),
        documentTitleIndex: documentTitleIndex,
    }
}

func (db *DocumentBuilder) CreateAndRegisterDocument(tx Transaction, compositePartID int) (Document, *OpFailedError) {
    id, err := db.idPool.GetID(tx)
    if err != nil {
        return nil, err
    }
    title := "Composite Part #" + strconv.Itoa(compositePartID)
    text := createText(internal.DocumentSize, "I am the documentation for composite part #" + strconv.Itoa(compositePartID) + "\n")
    document := doFactory.CreateDocument(tx, id, title, text)
    db.documentTitleIndex.Put(tx, title, document)
    return document, nil
}

func (db *DocumentBuilder) UnregisterAndRecycleDocument(tx Transaction, document Document) {
    document.SetPart(tx, nil)
    db.documentTitleIndex.Remove(tx, document.GetTitle(tx))
    db.idPool.PutUnusedID(tx, document.GetDocumentId(tx))
}
