package model

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD")).
			Padding(0, 1)

	noteListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)

	noteStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedNoteStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#874BFD")).
				Bold(true)
)

func (m Model) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Render(m.err.Error())
	}

	title := titleStyle.Render("Notes")

	var renderedNotes []string
	for _, note := range m.notes {
		renderedNotes = append(renderedNotes, noteStyle.Render(note.Name))
	}

	noteList := noteListStyle.Width(m.display.width / 2).Render(lipgloss.JoinVertical(lipgloss.Left, renderedNotes...))

	content := lipgloss.JoinVertical(lipgloss.Left, title, noteList)

	return appStyle.Width(m.display.width).Height(m.display.height).Render(content)
}
