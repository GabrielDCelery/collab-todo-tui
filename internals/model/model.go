package model

import (
	"context"
	"log"

	"github.com/GabrielDCelery/collab-todo-tui/internals/commands"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelConfig struct {
	dir               string
	allowedExtensions []string
}

type Subscriptions struct {
	noteCreatedChan chan commands.NoteCreated
	noteRemovedChan chan commands.NoteRemoved
}

type Model struct {
	ctx    context.Context
	err    error
	config ModelConfig
	notes  map[string]commands.Note
	subs   Subscriptions
}

func NewModel(ctx context.Context, dir string, allowedExtensions []string) Model {
	return Model{
		ctx: ctx,
		err: nil,
		config: ModelConfig{
			dir:               dir,
			allowedExtensions: allowedExtensions,
		},
		notes: map[string]commands.Note{},
		subs: Subscriptions{
			noteCreatedChan: make(chan commands.NoteCreated),
			noteRemovedChan: make(chan commands.NoteRemoved),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		commands.RegisterNotes(m.config.dir, m.config.allowedExtensions),
		commands.WaitForNoteCreated(m.subs.noteCreatedChan),
		commands.WaitForNoteRemoved(m.subs.noteRemovedChan),
	)
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := tea.NewProgram(NewModel(ctx, "/home/gzeller/github/collab-todo-tui", []string{".md"}))
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
