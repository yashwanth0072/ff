package actions

import (
	"os/exec"
	"syscall"
)

// HasMimeopen reports whether mimeopen (perl-file-mimeinfo) is installed.
func HasMimeopen() bool {
	_, err := exec.LookPath("mimeopen")
	return err == nil
}

// MimeOpenCmd returns an *exec.Cmd that runs `mimeopen -a path`, which
// always prompts with a numbered list of every registered handler for
// the file's mimetype instead of silently deferring to the default.
// The -a flag is required — bare `mimeopen` behaves exactly like
// xdg-open and asks nothing. Meant to be run via tea.ExecProcess since
// it reads the chosen number from stdin.
func MimeOpenCmd(path string) *exec.Cmd {
	return exec.Command("mimeopen", "-a", path)
}

// OpenDetached launches xdg-open on path without blocking or attaching to
// the current terminal, so the TUI keeps running underneath it. Works for
// both files (opens in the default app) and directories (opens in the
// default file manager) on any XDG-compliant desktop, Hyprland included.
func OpenDetached(path string) error {
	cmd := exec.Command("xdg-open", path)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	// Detach into its own session so it survives ff exiting and doesn't
	// receive signals meant for the TUI process group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return cmd.Start()
}
