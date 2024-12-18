package locks

import (
	. "github.com/NotARealMike/gvstm/stm"
	. "github.com/NotARealMike/gvstm/stmbench7/interfaces"
	"strings"
)

type documentImpl struct {
	id    int
	title string
	text  string
	part  CompositePart
}

func newDocumentImpl(tx Transaction, id int, title, text string) Document {
	return &documentImpl{
		id:    id,
		title: title,
		text:  text,
		part:  nil,
	}
}

func (d *documentImpl) SetPart(tx Transaction, part CompositePart) {
	d.part = part
}

func (d *documentImpl) GetCompositePart(tx Transaction) CompositePart {
	return d.part
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
	t := d.text
	for _, r := range t {
		if r == symbol {
			occurrences++
		}
	}
	return occurrences
}

func (d *documentImpl) ReplaceText(tx Transaction, from string, to string) int {
	t := d.text
	if !strings.HasPrefix(t, from) {
		return 0
	}
	d.text = strings.Replace(t, from, to, 1)
	return 1
}

func (d *documentImpl) TextBeginsWith(tx Transaction, prefix string) bool {
	return strings.HasPrefix(d.text, prefix)
}

func (d *documentImpl) GetText(tx Transaction) string {
	return d.text
}
