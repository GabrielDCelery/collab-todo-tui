package model

import (
	"context"
	"fmt"
	"log"

	"github.com/GabrielDCelery/collab-todo-tui/internals/notes"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	err          error
	notesWatcher *notes.NotesWatcher
	ctx          context.Context
}

func NewModel(ctx context.Context) Model {
	ti := textinput.New()
	ti.Placeholder = "Add new item"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		notesWatcher: notes.NewNotesWatcher("/home/gzeller/github/collab-todo-tui", []string{".md"}),
		err:          nil,
		ctx:          ctx,
	}
}

func (m Model) Init() tea.Cmd {
	err := m.notesWatcher.Start(m.ctx)
	if err != nil {
		return tea.Quit
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m Model) View() string {
	s := ""

	for _, note := range m.notesWatcher.Notes {
		s += fmt.Sprintf("%s\n", note.Name)
	}

	return s
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := tea.NewProgram(NewModel(ctx))
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
