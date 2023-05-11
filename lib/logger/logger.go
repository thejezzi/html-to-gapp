package logger

import (
	"fmt"
)

const (
	FG_RED     = 31
	FG_GREEN   = 32
	FG_YELLOW  = 33
	FG_BLUE    = 34
	FG_MAGENTA = 35
	FG_CYAN    = 36
	FG_WHITE   = 37

	FG_LIGHT_BLUE    = 94
	FG_LIGHT_GREEN   = 92
	FG_LIGHT_RED     = 91
	FG_LIGHT_YELLOW  = 93
	FG_LIGHT_MAGENTA = 95
	FG_LIGHT_CYAN    = 96

	BG_RED     = 41
	BG_GREEN   = 42
	BG_YELLOW  = 43
	BG_BLUE    = 44
	BG_MAGENTA = 45
	BG_CYAN    = 46
	BG_WHITE   = 47

	RESET = 0
)

func colorize(color int, message string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, message)
}

// func reset() string {
//   return colorize(RESET, "")
// }

func Info(message string) {
	outStr := "[" + colorize(FG_LIGHT_BLUE, "INFO") + "] " + message
	fmt.Println(outStr)
}

func Warn(message string) {
	outStr := "[" + colorize(FG_LIGHT_YELLOW, "WARN") + "] " + message
	fmt.Println(outStr)
}

func Error(message string) {
	outStr := "[" + colorize(FG_LIGHT_RED, "ERROR") + "] " + message
	fmt.Println(outStr)
}

func Debug(message string) {
	outStr := "[" + colorize(FG_LIGHT_MAGENTA, "DEBUG") + "] " + message
	fmt.Println(outStr)
}
