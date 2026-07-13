package actions

import "os/exec"

// UpdateDBCmd returns an *exec.Cmd that refreshes the plocate index.
// Run via tea.ExecProcess so sudo can prompt for a password on the
// real terminal.
func UpdateDBCmd() *exec.Cmd {
	return exec.Command("sudo", "updatedb")
}
