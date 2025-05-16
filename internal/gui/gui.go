package gui

import (
	"os/exec"
)

func KillGUI() error {
	exec.Command("systemctl", "isolate", "multi-user.target").Run()
	exec.Command("pkill", "-f", "Hyprland").Run()
	exec.Command("pkill", "-f", "gnome-shell").Run()
	return nil
}