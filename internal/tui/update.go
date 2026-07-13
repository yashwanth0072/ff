package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/yashwanth/ff/internal/actions"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp, tea.KeyCtrlK:
			m.moveCursor(-1)
			return m, nil

		case tea.KeyDown, tea.KeyCtrlJ:
			m.moveCursor(1)
			return m, nil

		case tea.KeyEnter:
			if path, ok := m.selected(); ok {
				if actions.HasMimeopen() {
					// mimeopen -a prompts interactively, so it needs the
					// terminal handed to it — same path as Ctrl+E.
					c := actions.MimeOpenCmd(path)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return execFinishedMsg{action: "open", err: err}
					})
				}
				// mimeopen not installed — fall back to silent xdg-open.
				if err := actions.OpenDetached(path); err != nil {
					m.statusMsg = "open failed: " + err.Error()
				} else {
					m.statusMsg = "opened " + path
				}
				return m, clearStatusCmd(m.generation)
			}
			return m, nil

		case tea.KeyCtrlE:
			if path, ok := m.selected(); ok {
				c := actions.EditorCmd(path)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return execFinishedMsg{action: "edit", err: err}
				})
			}
			return m, nil

		case tea.KeyCtrlY:
			if path, ok := m.selected(); ok {
				if err := actions.CopyPath(path); err != nil {
					m.statusMsg = "copy failed: " + err.Error()
				} else {
					m.statusMsg = "copied: " + path
				}
				return m, clearStatusCmd(m.generation)
			}
			return m, nil

		case tea.KeyCtrlU:
			c := actions.UpdateDBCmd()
			m.statusMsg = "refreshing plocate index (sudo)..."
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return execFinishedMsg{action: "updatedb", err: err}
			})
		}

	case debounceMsg:
		if msg.gen != m.generation {
			return m, nil // stale, a newer keystroke already superseded this
		}
		m.searching = true
		return m, searchCmd(m.input.Value(), m.generation)

	case resultsMsg:
		if msg.gen != m.generation {
			return m, nil // stale results from an older keystroke
		}
		m.searching = false
		m.err = msg.err
		m.results = msg.results
		m.cursor = 0
		m.offset = 0
		return m, nil

	case execFinishedMsg:
		if msg.err != nil {
			m.statusMsg = msg.action + " failed: " + msg.err.Error()
		} else {
			m.statusMsg = msg.action + " done"
		}
		return m, clearStatusCmd(m.generation)

	case statusClearMsg:
		if msg.gen == m.generation {
			m.statusMsg = ""
		}
		return m, nil
	}

	prevValue := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if m.input.Value() != prevValue {
		m.generation++
		m.statusMsg = ""
		if m.input.Value() == "" {
			m.results = nil
			m.searching = false
			return m, cmd
		}
		return m, tea.Batch(cmd, debounceCmd(m.generation))
	}

	return m, cmd
}

func (m *Model) moveCursor(delta int) {
	if len(m.results) == 0 {
		return
	}
	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor > len(m.results)-1 {
		m.cursor = len(m.results) - 1
	}

	visible := m.listHeight()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m Model) selected() (string, bool) {
	if m.cursor < 0 || m.cursor >= len(m.results) {
		return "", false
	}
	return m.results[m.cursor].Path, true
}

func (m Model) listHeight() int {
	h := m.height - 6 // room for input box, status line, help line, borders
	if h < 3 {
		h = 3
	}
	return h
}
