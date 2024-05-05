package interfaces

import (
	. "github.com/NotARealMike/gvstm/stm"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"strconv"
)

type documentBuilder interface {
	createAndRegisterDocument(tx Transaction, compositePartID int) (Document, OpFailedError)
	unregisterAndRecycleDocument(tx Transaction, document Document)
}

type documentBuilderImpl struct {
	idPool             IDPool
	documentTitleIndex Index
}

func newDocumentBuilder(tx Transaction, documentTitleIndex Index) documentBuilder {
	return &documentBuilderImpl{
		idPool:             BEFactory.CreateIDPool(tx, internal.MaxCompParts),
		documentTitleIndex: documentTitleIndex,
	}
}

func (db *documentBuilderImpl) createAndRegisterDocument(tx Transaction, compositePartID int) (Document, OpFailedError) {
	id, err := db.idPool.GetID(tx)
	if err != nil {
		return nil, err
	}
	title := "Composite Part #" + strconv.Itoa(compositePartID)
	text := createText(internal.DocumentSize, "I am the documentation for composite part #"+strconv.Itoa(compositePartID)+"\n")
	document := DOFactory.CreateDocument(tx, id, title, text)
	db.documentTitleIndex.Put(tx, title, document)
	return document, nil
}

func (db *documentBuilderImpl) unregisterAndRecycleDocument(tx Transaction, document Document) {
	document.SetPart(tx, nil)
	db.documentTitleIndex.Remove(tx, document.GetTitle(tx))
	db.idPool.PutUnusedID(tx, document.GetDocumentId(tx))
}
