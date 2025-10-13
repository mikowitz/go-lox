package lox

import (
	"io"
	"os"
)

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	err = f()
	os.Stdout = orig
	w.Close()
	out, readErr := io.ReadAll(r)
	if readErr != nil {
		return "", readErr
	}
	return string(out), err
}
