package interactive

import (
	"bufio"
	"fmt"
	"os"

	"sysfail/internal/signals"
)

func Run(lines []string) {
	signals.Ignore()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\033[2J\033[H")
		for _, line := range lines {
			fmt.Println(line)
		}
		fmt.Print("Press any key to continue...")
		_, _ = reader.ReadByte()
	}
}
