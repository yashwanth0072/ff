package actions

import (
	"errors"
	"os/exec"
	"strings"
)

// CopyPath copies path to the system clipboard. Tries wl-copy first
// (Hyprland/Wayland), then falls back to xclip/xsel for X11 sessions.
func CopyPath(path string) error {
	candidates := [][]string{
		{"wl-copy"},
		{"xclip", "-selection", "clipboard"},
		{"xsel", "--clipboard", "--input"},
	}

	var lastErr error
	for _, c := range candidates {
		bin := c[0]
		if _, err := exec.LookPath(bin); err != nil {
			continue
		}
		cmd := exec.Command(bin, c[1:]...)
		cmd.Stdin = strings.NewReader(path)
		if err := cmd.Run(); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	return errors.New("no clipboard tool found (install wl-clipboard: sudo pacman -S wl-clipboard)")
}
