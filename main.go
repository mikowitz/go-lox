package main

import (
	"os"

	"github.com/mikowitz/go-lox/internal/lox"
)

func main() {
	loxRuntime := &lox.Lox{}
	os.Exit(loxRuntime.Run(os.Args[1:]))
}
