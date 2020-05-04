package gvstm

import (
	"gvstm/gvstm"
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
	"strings"
)

type documentImpl struct {
	id    int
	title string
	text  TVar
	part  TVar
}

func newDocumentImpl(tx Transaction, id int, title, text string) Document {
	return &documentImpl{
		id:    id,
		title: title,
		text:  gvstm.CreateTVar(text),
		part:  gvstm.CreateTVar(nil),
	}
}

func (d *documentImpl) SetPart(tx Transaction, part CompositePart) {
	tx.Store(d.part, part)
}

func (d *documentImpl) GetCompositePart(tx Transaction) CompositePart {
	return tx.Load(d.part).(CompositePart)
}

func (d *documentImpl) GetDocumentId(tx Transaction) int {
	return d.id
}

func (d *documentImpl) GetTitle(tx Transaction) string {
	return d.title
}

func (d *documentImpl) NullOperation(tx Transaction) {}

func (d *documentImpl) SearchText(tx Transaction, symbol rune) int {
	occurrences := 0
	t := tx.Load(d.text).(string)
	for _, r := range t {
		if r == symbol {
			occurrences++
		}
	}
	return occurrences
}

func (d *documentImpl) ReplaceText(tx Transaction, from string, to string) int {
	t := tx.Load(d.text).(string)
	if !strings.HasPrefix(t, from) {
		return 0
	}
	tx.Store(d.text, strings.Replace(t, from, to, 1))
	return 1
}

func (d *documentImpl) TextBeginsWith(tx Transaction, prefix string) bool {
	return strings.HasPrefix(tx.Load(d.text).(string), prefix)
}

func (d *documentImpl) GetText(tx Transaction) string {
	return tx.Load(d.text).(string)
}
