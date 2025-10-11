// Package lox
package lox

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ExitUsage   = 64
	ExitDataErr = 65
)

type Lox struct {
	hadError bool
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) Run(args []string) int {
	var exitStatus int
	switch len(args) {
	case 0:
		exitStatus = l.runPrompt()
	case 1:
		exitStatus = l.runFile(args[0])
	default:
		fmt.Println("Usage: go-lox [script]")
		exitStatus = ExitUsage
	}
	return exitStatus
}

func (l *Lox) runPrompt() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			break
		}
		l.run(strings.TrimSpace(input))
		l.hadError = false
	}
	return 0
}

func (l *Lox) runFile(filepath string) int {
	if f, err := os.ReadFile(filepath); err == nil {
		l.run(string(f))
	}
	if l.hadError {
		// os.Exit(ExitDataErr)
		return ExitDataErr
	}
	return 0
}

func (l *Lox) run(input string) {
	// fmt.Printf("running:\n\n%s\n", input)
}
