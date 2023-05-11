package cli

import (
	"fmt"
)

func PrintHelp() {
	fmt.Println("Usage: htmltogo <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  convert <html file> <gapp file>")
	fmt.Println("    Convert an html file to a gapp file")
	fmt.Println("")
	fmt.Println("  help")
	fmt.Println("    Print this help message")
	fmt.Println("")
}
