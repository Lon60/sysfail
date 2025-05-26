package modes

import (
	"fmt"
	"time"

	"sysfail/internal/gui"
	"sysfail/internal/interactive"
	panicpkg "sysfail/internal/panic"
	"sysfail/internal/signals"
	"sysfail/internal/system"
)

type Config struct {
	Mode      string
	Duration  time.Duration
	PanicType string
}

func Execute(config *Config) error {
	switch config.Mode {
	case "console":
		return executeConsoleMode(config)
	default:
		return fmt.Errorf("unknown mode '%s'. Available modes: console", config.Mode)
	}
}

func executeConsoleMode(config *Config) error {
	sysInfo := system.DetectSystem()

	panicLines, exists := panicpkg.GetPanic(config.PanicType, sysInfo)
	if !exists {
		return fmt.Errorf("unknown panic type '%s'. Available: panic, oops, segfault", config.PanicType)
	}

	signals.Ignore()

	if err := gui.KillGUI(); err != nil {
	}

	fmt.Print("\033[2J\033[H")
	time.Sleep(500 * time.Millisecond)

	if config.Duration > 0 {
		go func() {
			time.Sleep(config.Duration)
			fmt.Println("\n[TIMEOUT] Emergency session terminated")
			interactive.GracefulExit()
		}()
	}

	interactive.Run(panicLines, sysInfo)

	return nil
}
