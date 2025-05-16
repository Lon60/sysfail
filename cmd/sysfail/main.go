package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"sysfail/internal/gui"
	"sysfail/internal/interactive"
	panicpkg "sysfail/internal/panic"
	"sysfail/internal/signals"
)

func main() {

	errorType := flag.String("error", "panic", "which fake panic to display (panic|oops)")
	interactiveMode := flag.Bool("interactive", false, "go interactive and block exit")
	killGui := flag.Bool("kill-gui", false, "tear down graphical interface (requires sudo)")
	flag.Parse()


	lines, ok := panicpkg.Panics[*errorType]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown error type: %s\n", *errorType)
		os.Exit(2)
	}


	if *killGui {
		fmt.Println("[sysfail] killing GUI…")
		if err := gui.KillGUI(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to kill GUI: %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}

	signals.Ignore()

	if *interactiveMode {
		interactive.Run(lines)
	} else {
		for _, line := range lines {
			if isHeader(line) {
				fmt.Printf("\033[1;31m%s\033[0m\n", line)
			} else {
				fmt.Println(line)
			}
			time.Sleep(300 * time.Millisecond)
		}
		os.Exit(1)
	}
}

func isHeader(s string) bool {
	return len(s)>0 && (strings.Contains(strings.ToLower(s), "panic") || strings.Contains(strings.ToLower(s), "oops"))
}