package signals

import (
	"os/signal"
	"syscall"
)

func Ignore() {
	signal.Ignore(syscall.SIGINT, syscall.SIGTSTP)
}