package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/yashwanth/ff/internal/search"
)

const debounceDelay = 120 * time.Millisecond

// Model is the ff application state.
type Model struct {
	input   textinput.Model
	results []search.Result
	cursor  int
	offset  int // scroll offset for the visible window

	width  int
	height int

	generation int // bumped on every keystroke; stale search results are dropped
	searching  bool
	err        error
	statusMsg  string
}

func New(initialQuery string) Model {
	ti := textinput.New()
	ti.Placeholder = "type to search the whole disk..."
	ti.Prompt = "🔎 "
	ti.Focus()
	ti.CharLimit = 256
	ti.SetValue(initialQuery)

	m := Model{
		input:  ti,
		width:  80,
		height: 24,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	if m.input.Value() != "" {
		return searchCmd(m.input.Value(), m.generation)
	}
	return textinput.Blink
}

// --- messages -------------------------------------------------------

type debounceMsg struct{ gen int }

type resultsMsg struct {
	gen     int
	results []search.Result
	err     error
}

type execFinishedMsg struct {
	action string
	err    error
}

type statusClearMsg struct{ gen int }

// --- commands ---------------------------------------------------------

func debounceCmd(gen int) tea.Cmd {
	return tea.Tick(debounceDelay, func(time.Time) tea.Msg {
		return debounceMsg{gen: gen}
	})
}

func searchCmd(query string, gen int) tea.Cmd {
	return func() tea.Msg {
		candidates, err := search.Candidates(query)
		if err != nil {
			return resultsMsg{gen: gen, err: err}
		}
		ranked := search.Rank(candidates, query)
		return resultsMsg{gen: gen, results: ranked}
	}
}

func clearStatusCmd(gen int) tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return statusClearMsg{gen: gen}
	})
}
