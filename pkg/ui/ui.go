package ui

import (
	"fmt"
	"os"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// Banner prints the stylish ASCII art header.
func Banner(version string) {
	fmt.Println()
	fmt.Printf("%s%s  ____      ____ _                  %s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s / ___|___ / ___| | ___  __ _ _ __  %s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s| |  _/ _ \\ |   | |/ _ \\/ _` | '_ \\ %s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s| |_| | (_) | |___| |  __/ (_| | | | |%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s \\____|\\___/ \\____|_|\\___|\\__,_|_| |_|%s  %s%sv%s%s\n", Bold, Cyan, Reset, Dim, White, version, Reset)
	fmt.Printf("%s  Go Clean Architecture Scaffolding Tool%s\n", Dim, Reset)
	fmt.Println()
}

// Success prints a green checkmark message.
func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s✓%s %s\n", Bold, Green, Reset, msg)
}

// Info prints a cyan info message.
func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%sℹ%s %s\n", Bold, Cyan, Reset, msg)
}

// Warn prints a yellow warning message.
func Warn(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s!%s %s\n", Bold, Yellow, Reset, msg)
}

// Step prints a step indicator.
func Step(step int, total int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s[%d/%d]%s %s...\n", Cyan, step, total, Reset, msg)
}

// Error prints an error message to stderr.
func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "  %s%s✖ Error:%s %s\n", Bold, Red, Reset, msg)
}

// Header prints a section header.
func Header(title string) {
	fmt.Printf("\n%s%s%s%s\n", Bold, White, title, Reset)
}
