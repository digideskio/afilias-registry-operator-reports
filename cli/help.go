package cli

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"os"
)

func HelpBanner(command string) {
	color.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s %s\n", os.Args[0], command))
}
