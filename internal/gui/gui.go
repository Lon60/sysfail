package gui

import (
	"errors"
	"os/exec"
)

func KillGUI() error {

	commands := [][]string{
		{"pkill", "-f", "gnome-session"},
		{"pkill", "-f", "kde-plasma"},
		{"pkill", "-f", "xfce4-session"},
		{"pkill", "-f", "X"},
		{"systemctl", "stop", "gdm"},
		{"systemctl", "stop", "lightdm"},
		{"systemctl", "stop", "sddm"},
	}

	var lastErr error
	for _, cmd := range commands {
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err == nil {
			return nil // Success!
		}
		lastErr = err
	}

	if lastErr != nil {
		return errors.New("failed to kill GUI - try running with sudo")
	}

	return nil
}
