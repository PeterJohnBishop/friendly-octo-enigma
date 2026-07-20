// Package tui provides the terminal user interface
package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/peterjohnbishop/friendly-octo-enigma/models"
)

type (
	headerMsg map[string][]string
	bodyMsg   map[string]models.PropertyDetail
)

type listItem struct {
	key      string
	val      any
	display  string
	selected bool
}

const maxVisible = 10

type sessionState int

const (
	stateSelecting sessionState = iota // 0
	stateNextStage                     // 1
)

type uiModel struct {
	headersChan <-chan map[string][]string
	bodyChan    <-chan map[string]models.PropertyDetail

	headers []listItem
	body    []listItem

	focus        int
	headerCursor int
	headerOffset int
	bodyCursor   int
	bodyOffset   int

	state sessionState
	args  []any

	// Add these for dynamic resizing:
	windowWidth  int
	windowHeight int
}

func (m uiModel) Init() tea.Cmd {
	return tea.Batch(
		listenHeaders(m.headersChan),
		listenBody(m.bodyChan),
	)
}

func Start(hChan <-chan map[string][]string, bChan <-chan map[string]models.PropertyDetail) error {
	m := uiModel{
		headersChan: hChan,
		bodyChan:    bChan,
	}
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
