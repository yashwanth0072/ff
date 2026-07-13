// Command ff is a stupid-fast fuzzy file finder for the whole filesystem.
// It uses plocate for instant candidate retrieval and an in-memory fuzzy
// matcher (ported from the same approach as pkgmngr) for ranking.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yashwanth/ff/internal/tui"
)

func main() {
	initialQuery := strings.Join(os.Args[1:], " ")

	m := tui.New(initialQuery)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ff: "+err.Error())
		os.Exit(1)
	}
}
