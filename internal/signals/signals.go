package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func Ignore() {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Terminal kill
		syscall.SIGQUIT, // Ctrl+\
	)

	go func() {
		for range c {
			// Just ignore the signal - don't exit
		}
	}()
}
