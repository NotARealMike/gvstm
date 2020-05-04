package gvstm

import (
	"gvstm/gvstm"
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
	"strings"
)

type manualImpl struct {
	id     int
	title  string
	text   TVar
	module TVar
}

func newManualImpl(tx Transaction, id int, title, text string) Manual {
	return &manualImpl{
		id:     id,
		title:  title,
		text:   gvstm.CreateTVar(text),
		module: gvstm.CreateTVar(nil),
	}
}

func (m *manualImpl) GetId(tx Transaction) int {
	return m.id
}

func (m *manualImpl) GetTitle(tx Transaction) string {
	return m.title
}

func (m *manualImpl) GetText(tx Transaction) string {
	return tx.Load(m.text).(string)
}

func (m *manualImpl) GetModule(tx Transaction) Module {
	return tx.Load(m.module).(Module)
}

func (m *manualImpl) SetModule(tx Transaction, module Module) {
	tx.Store(m.module, module)
}

func (m *manualImpl) CountOccurences(tx Transaction, ch rune) int {
	occurrences := 0
	t := tx.Load(m.text).(string)
	for _, r := range t {
		if r == ch {
			occurrences++
		}
	}
	return occurrences
}

func (m *manualImpl) CheckFirstLastCharTheSame(tx Transaction) int {
	t := tx.Load(m.text).(string)
	if t[0] == t[len(t)-1] {
		return 1
	}
	return 0
}

func (m *manualImpl) StartsWith(tx Transaction, ch rune) bool {
	return strings.HasPrefix(tx.Load(m.text).(string), string(ch))
}

func (m *manualImpl) ReplaceChar(tx Transaction, from rune, to rune) int {
	tx.Store(m.text, strings.ReplaceAll(tx.Load(m.text).(string), string(from), string(to)))
	return m.CountOccurences(tx, to)
}
