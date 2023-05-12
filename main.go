package main

import (
	"fmt"
	"os"

	"github.com/thejezzi/html-to-gapp/lib/cli"
	"github.com/thejezzi/html-to-gapp/lib/htg"
	"github.com/thejezzi/html-to-gapp/lib/logger"
)

func main() {

	// TestRun()
	htg.TestRun()
	os.Exit(0)

	// Simple cli boilerplate
	if len(os.Args) < 2 {
		fmt.Println("Usage: htmltogo <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "convert":
		if len(os.Args) < 4 {
			logger.Error("Usage: htmltogo convert <html file> <gapp file>")
			cli.PrintHelp()
			os.Exit(1)
		}
		htg.Convert(os.Args[2], os.Args[3])
		logger.Error("Not implemented yet")
	default:
		cli.PrintHelp()
		os.Exit(1)
	}

}
