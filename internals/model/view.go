package model

import (
	"fmt"
)

func (m Model) View() string {
	s := ""

	if m.err != nil {
		s += m.err.Error()
		return s
	}

	for _, note := range m.notes {
		s += fmt.Sprintf("%s\n", note.Name)
	}

	return s

}
