package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("229")).
			Bold(true)

	matchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ff") + dimStyle.Render(" — fuzzy file finder") + "\n")
	b.WriteString(borderStyle.Width(m.width-4).Render(m.input.View()) + "\n")

	if m.err != nil {
		b.WriteString(errStyle.Render("error: "+m.err.Error()) + "\n")
	}

	b.WriteString(m.renderList())

	status := m.statusMsg
	if status == "" && m.searching {
		status = "searching..."
	}
	if status == "" && m.input.Value() != "" {
		status = countLabel(len(m.results))
	}
	b.WriteString(statusStyle.Render(status) + "\n")

	b.WriteString(helpStyle.Render("↑/↓ navigate · enter open · ctrl+e edit · ctrl+y copy path · ctrl+u refresh index · esc quit"))

	return b.String()
}

func countLabel(n int) string {
	if n == 1 {
		return "1 match"
	}
	if n == 0 {
		return "no matches"
	}
	if n >= plocateResultCap {
		return itoa(plocateResultCap) + "+ matches — keep typing to narrow it down"
	}
	return itoa(n) + " matches"
}

// kept small and dependency-free instead of pulling in fmt.Sprintf everywhere
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// Mirrors the -l limit passed to plocate in internal/search/plocate.go.
const plocateResultCap = 1500

func (m Model) renderList() string {
	if len(m.results) == 0 {
		if m.input.Value() == "" {
			return dimStyle.Render("  start typing to search your whole filesystem\n")
		}
		if !m.searching {
			return dimStyle.Render("  nothing found\n")
		}
		return "\n"
	}

	visible := m.listHeight()
	end := m.offset + visible
	if end > len(m.results) {
		end = len(m.results)
	}

	var b strings.Builder
	for i := m.offset; i < end; i++ {
		r := m.results[i]
		line := renderPath(r.Path, r.MatchedIndexes)

		icon := "📄"
		if info, err := os.Stat(r.Path); err == nil && info.IsDir() {
			icon = "📁"
		}

		row := "  " + icon + " " + line
		if i == m.cursor {
			row = selectedStyle.Render("▶ " + icon + " " + line)
		}
		b.WriteString(row + "\n")
	}
	return b.String()
}

// renderPath bolds the characters at matchedIndexes within path.
func renderPath(path string, matchedIndexes []int) string {
	if len(matchedIndexes) == 0 {
		return path
	}
	matchSet := make(map[int]bool, len(matchedIndexes))
	for _, idx := range matchedIndexes {
		matchSet[idx] = true
	}

	var b strings.Builder
	for i, r := range []rune(path) {
		if matchSet[i] {
			b.WriteString(matchStyle.Render(string(r)))
		} else {
			b.WriteString(string(r))
		}
	}
	return b.String()
}
