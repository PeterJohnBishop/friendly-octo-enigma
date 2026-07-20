package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	pink           = lipgloss.Color("212") // A vibrant pink
	activeCursor   = lipgloss.NewStyle().Foreground(pink).Bold(true)
	selectedItem   = lipgloss.NewStyle().Foreground(pink)
	unselectedItem = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	activeBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Padding(0, 1).
			Width(50)

	inactiveBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(50)
)

func (m uiModel) View() tea.View {
	var b strings.Builder
	b.WriteString("\n  Listening for Webhooks... (Press 'Tab' to switch, 'Space' to select, 'Enter' to confirm)\n\n")

	maxVis := calculateMaxVisible(m.windowHeight)

	boxWidth := m.windowWidth - 4
	if boxWidth < 50 {
		boxWidth = 50
	}
	activeBox = activeBox.Width(boxWidth)
	inactiveBox = inactiveBox.Width(boxWidth)

	headerTitle := "Headers"
	headerList := renderList(m.headers, m.headerCursor, m.headerOffset, maxVis, m.focus == 0)
	if m.focus == 0 {
		b.WriteString(activeBox.Render(headerTitle+"\n\n"+headerList) + "\n")
	} else {
		b.WriteString(inactiveBox.Render(headerTitle+"\n\n"+headerList) + "\n")
	}

	bodyTitle := "Body Properties"
	bodyList := renderList(m.body, m.bodyCursor, m.bodyOffset, maxVis, m.focus == 1)
	if m.focus == 1 {
		b.WriteString(activeBox.Render(bodyTitle+"\n\n"+bodyList) + "\n")
	} else {
		b.WriteString(inactiveBox.Render(bodyTitle+"\n\n"+bodyList) + "\n")
	}

	return tea.NewView(b.String())
}
