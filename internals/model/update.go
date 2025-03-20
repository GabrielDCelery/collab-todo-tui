package model

import (
	"github.com/GabrielDCelery/collab-todo-tui/internals/commands"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.display.width = msg.Width
		m.display.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case commands.NotesRegistered:
		m.notes = msg
		return m, commands.ListenForNoteChanges(
			m.ctx,
			m.config.dir,
			m.config.allowedExtensions,
			m.subs.noteCreatedChan,
			m.subs.noteRemovedChan,
		)

	case commands.NoteCreated:
		m.notes[msg.Path] = msg.Note
		return m, commands.WaitForNoteCreated(m.subs.noteCreatedChan)

	case commands.NoteRemoved:
		delete(m.notes, msg.Path)
		return m, commands.WaitForNoteRemoved(m.subs.noteRemovedChan)

	case error:
		m.err = msg
		return m, nil

	}

	return m, cmd
}
