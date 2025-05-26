package main

import (
	"flag"
	"fmt"
	"os"
	"sysfail/internal/modes"
)

const version = "1.0.0"

func main() {
	var (
		mode      = flag.String("m", "", "execution mode (required)")
		duration  = flag.Duration("d", 0, "auto-exit after duration (0 = never)")
		panicType = flag.String("t", "panic", "panic type (panic|oops|segfault)")
		help      = flag.Bool("h", false, "show help")
	)

	flag.StringVar(mode, "mode", "", "execution mode (required)")
	flag.DurationVar(duration, "duration", 0, "auto-exit after duration (0 = never)")
	flag.StringVar(panicType, "type", "panic", "panic type (panic|oops|segfault)")
	flag.BoolVar(help, "help", false, "show help")

	flag.Parse()

	if *help || *mode == "" {
		printHelp()
		if *mode == "" {
			os.Exit(1)
		}
		return
	}

	config := &modes.Config{
		Mode:      *mode,
		Duration:  *duration,
		PanicType: *panicType,
	}

	if err := modes.Execute(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`sysfail v%s - Linux System Failure Simulator

USAGE:
    sysfail -m <mode> [options]

MODES:
    console     Interactive fake recovery console (kills GUI, shows panic, traps user)

OPTIONS:
    -m, --mode      Execution mode (required)
    -d, --duration  Auto-exit after duration (e.g. 30s, 5m) [default: never]
    -t, --type      Panic type: panic, oops, segfault [default: panic]
    -h, --help      Show this help

EXAMPLES:
    sysfail -m console                    # Basic console mode
    sysfail -m console -d 2m              # Exit after 2 minutes
    sysfail -m console -t oops -d 30s     # Oops panic, exit after 30 seconds

WARNING: This is a prank tool. Use responsibly!
`, version)
}
