package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/peterjohnbishop/friendly-octo-enigma/models"
)

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case headerMsg:
		m.headers = buildListItems(msg, func(k string, v []string) string {
			return fmt.Sprintf("%s: %v", k, strings.Join(v, ", "))
		})
		m.headerCursor = 0
		m.headerOffset = 0
		return m, listenHeaders(m.headersChan)

	case bodyMsg:
		m.body = buildListItems(msg, func(k string, v models.PropertyDetail) string {
			formattedVal := formatPropertyValue(v)
			return fmt.Sprintf("%s: %s", k, formattedVal)
		})
		m.bodyCursor = 0
		m.bodyOffset = 0
		return m, listenBody(m.bodyChan)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.state == stateSelecting {
				m.args = m.GetArgs()
				if len(m.args) > 0 {
					m.state = stateNextStage
				}
			}

		case "tab":
			if m.state == stateSelecting {
				m.focus = (m.focus + 1) % 2
			}

		case " ", "space":
			if m.state == stateSelecting {
				if m.focus == 0 && len(m.headers) > 0 {
					m.headers[m.headerCursor].selected = !m.headers[m.headerCursor].selected
				} else if m.focus == 1 && len(m.body) > 0 {
					m.body[m.bodyCursor].selected = !m.body[m.bodyCursor].selected
				}
			}

		case "up", "k":
			if m.state == stateSelecting {
				// (Keep your existing 'up' offset/cursor logic here)
				if m.focus == 0 && m.headerCursor > 0 {
					m.headerCursor--
					if m.headerCursor < m.headerOffset {
						m.headerOffset--
					}
				} else if m.focus == 1 && m.bodyCursor > 0 {
					m.bodyCursor--
					if m.bodyCursor < m.bodyOffset {
						m.bodyOffset--
					}
				}
			}

		case "down", "j":
			if m.state == stateSelecting {
				maxVisible := calculateMaxVisible(m.windowHeight)

				if m.focus == 0 && m.headerCursor < len(m.headers)-1 {
					m.headerCursor++
					if m.headerCursor >= m.headerOffset+maxVisible {
						m.headerOffset++
					}
				} else if m.focus == 1 && m.bodyCursor < len(m.body)-1 {
					m.bodyCursor++
					if m.bodyCursor >= m.bodyOffset+maxVisible {
						m.bodyOffset++
					}
				}
			}
		}
	}

	return m, nil
}

func listenHeaders(ch <-chan map[string][]string) tea.Cmd {
	return func() tea.Msg { return headerMsg(<-ch) }
}

func listenBody(ch <-chan map[string]models.PropertyDetail) tea.Cmd {
	return func() tea.Msg { return bodyMsg(<-ch) }
}

func buildListItems[T any](data map[string]T, format func(string, T) string) []listItem {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var items []listItem
	for _, k := range keys {
		items = append(items, listItem{
			key:      k,
			val:      data[k],
			display:  format(k, data[k]),
			selected: false,
		})
	}
	return items
}

func renderList(items []listItem, cursor int, offset int, maxVisible int, isFocused bool) string {
	if len(items) == 0 {
		return unselectedItem.Render("  (Waiting for data...)")
	}

	var s strings.Builder

	if offset > 0 {
		s.WriteString(unselectedItem.Render("    ↑") + "\n")
	} else {
		s.WriteString("\n")
	}

	end := offset + maxVisible
	if end > len(items) {
		end = len(items)
	}

	for i := offset; i < end; i++ {
		item := items[i]

		cursorStr := "  "
		if isFocused && cursor == i {
			cursorStr = "> "
		}

		checkStr := "[ ] "
		if item.selected {
			checkStr = "[✓] "
		}

		line := fmt.Sprintf("%s%s%s", cursorStr, checkStr, item.display)

		if item.selected || (isFocused && cursor == i) {
			s.WriteString(selectedItem.Render(line) + "\n")
		} else {
			s.WriteString(unselectedItem.Render(line) + "\n")
		}
	}

	for i := len(items); i < maxVisible; i++ {
		s.WriteString("\n")
	}

	if end < len(items) {
		s.WriteString(unselectedItem.Render("    ↓") + "\n")
	} else {
		s.WriteString("\n")
	}

	return s.String()
}

func (m uiModel) GetArgs() []any {
	var args []any
	for _, h := range m.headers {
		if h.selected {
			args = append(args, h.val)
		}
	}
	for _, b := range m.body {
		if b.selected {
			args = append(args, b.val)
		}
	}
	return args
}

func calculateMaxVisible(windowHeight int) int {
	max := (windowHeight - 12) / 2
	if max < 3 {
		return 3
	}
	return max
}

func formatPropertyValue(v models.PropertyDetail) string {
	switch val := v.Value.(type) {

	case []models.PropertyDetail:
		var elements []string
		for _, item := range val {
			if item.InferredType == "string" {
				elements = append(elements, fmt.Sprintf(`"%v" [%v]`, item.Value, item.InferredType))
			} else {
				elements = append(elements, fmt.Sprintf("%v [%v]", item.Value, item.InferredType))
			}
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case []any:
		var elements []string
		for _, item := range val {
			switch inner := item.(type) {
			case string:
				elements = append(elements, fmt.Sprintf(`"%s" [string]`, inner))
			case models.PropertyDetail: // In case PropertyDetail structs are hiding in an 'any' slice
				if inner.InferredType == "string" {
					elements = append(elements, fmt.Sprintf(`"%v" [%v]`, inner.Value, inner.InferredType))
				} else {
					elements = append(elements, fmt.Sprintf("%v [%v]", inner.Value, inner.InferredType))
				}
			default:
				elements = append(elements, fmt.Sprintf("%v [%T]", item, item))
			}
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case []string:
		var elements []string
		for _, str := range val {
			elements = append(elements, fmt.Sprintf(`"%s" [string]`, str))
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case string:
		return fmt.Sprintf(`"%s" [%s]`, val, v.InferredType)

	default:
		return fmt.Sprintf("%v [%s]", val, v.InferredType)
	}
}
