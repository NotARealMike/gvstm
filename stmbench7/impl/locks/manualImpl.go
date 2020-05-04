package locks

import (
	. "gvstm/stm"
	. "gvstm/stmbench7/interfaces"
	"strings"
)

type manualImpl struct {
	id     int
	title  string
	text   string
	module Module
}

func newManualImpl(tx Transaction, id int, title, text string) Manual {
	return &manualImpl{
		id:     id,
		title:  title,
		text:   text,
		module: nil,
	}
}

func (m *manualImpl) GetId(tx Transaction) int {
	return m.id
}

func (m *manualImpl) GetTitle(tx Transaction) string {
	return m.title
}

func (m *manualImpl) GetText(tx Transaction) string {
	return m.text
}

func (m *manualImpl) GetModule(tx Transaction) Module {
	return m.module
}

func (m *manualImpl) SetModule(tx Transaction, module Module) {
	m.module = module
}

func (m *manualImpl) CountOccurences(tx Transaction, ch rune) int {
	occurrences := 0
	for _, r := range m.text {
		if r == ch {
			occurrences++
		}
	}
	return occurrences
}

func (m *manualImpl) CheckFirstLastCharTheSame(tx Transaction) int {
	if m.text[0] == m.text[len(m.text)-1] {
		return 1
	}
	return 0
}

func (m *manualImpl) StartsWith(tx Transaction, ch rune) bool {
	return strings.HasPrefix(m.text, string(ch))
}

func (m *manualImpl) ReplaceChar(tx Transaction, from rune, to rune) int {
	m.text = strings.ReplaceAll(m.text, string(from), string(to))
	return m.CountOccurences(tx, to)
}
