package actions

import (
	"os"
	"os/exec"
)

// EditorCmd returns an *exec.Cmd that opens path in $EDITOR, falling back
// to nvim then vi. This is meant to be run via tea.ExecProcess, which
// suspends the TUI and hands the terminal to the child process.
func EditorCmd(path string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if _, err := exec.LookPath("nvim"); err == nil {
			editor = "nvim"
		} else {
			editor = "vi"
		}
	}
	return exec.Command(editor, path)
}
