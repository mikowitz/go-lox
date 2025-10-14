package lox

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLox_error(t *testing.T) {
	l := &Lox{}

	require.False(t, l.hadError, "expected hadError to be false initially")

	l.error(1, "test error")

	assert.True(t, l.hadError, "expected hadError to be true after error")
}

func TestLox_report(t *testing.T) {
	tests := []struct {
		name    string
		line    int
		where   string
		message string
	}{
		{
			name:    "basic error",
			line:    1,
			where:   "",
			message: "Unexpected token",
		},
		{
			name:    "error with location",
			line:    5,
			where:   "at 'foo'",
			message: "Undefined variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lox{}

			require.False(t, l.hadError, "expected hadError to be false initially")

			l.report(tt.line, tt.where, tt.message)

			assert.True(t, l.hadError, "expected hadError to be true after report")
		})
	}
}

func TestLox_Run(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		fileContent        string
		expectHadError     bool
		expectedExitStatus int
	}{
		{
			name:               "no args - prompt mode",
			args:               []string{},
			expectHadError:     false,
			expectedExitStatus: 0,
		},
		{
			name:               "single file arg",
			args:               []string{"test.lox"},
			fileContent:        "true == false",
			expectHadError:     false,
			expectedExitStatus: 0,
		},
		{
			name:               "multiple single file args",
			args:               []string{"test.lox", "foo.lox"},
			fileContent:        "print \"hello\";",
			expectHadError:     true,
			expectedExitStatus: ExitUsage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lox{}

			// For file-based tests, create a temporary file
			if len(tt.args) == 1 {
				tmpDir := t.TempDir()
				tmpFile := filepath.Join(tmpDir, tt.args[0])

				err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644)
				require.NoError(t, err, "failed to create temp file")

				tt.args[0] = tmpFile

				exitStatus := l.Run(tt.args)

				assert.Equal(t, tt.expectHadError, l.hadError)
				assert.Equal(t, tt.expectedExitStatus, exitStatus)
			}
		})
	}
}

func TestLox_runFile(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectHadError bool
	}{
		{
			name:           "empty file",
			content:        "",
			expectHadError: false,
		},
		{
			name:           "simple content",
			content:        "3 + 4",
			expectHadError: false,
		},
		{
			name:           "multiline content",
			content:        "3 * 2 -\n 17",
			expectHadError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lox{}

			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.lox")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			require.NoError(t, err, "failed to create temp file")

			l.runFile(tmpFile)

			assert.Equal(t, tt.expectHadError, l.hadError)
		})
	}
}

func TestLox_run(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "simple statement",
			input: "3 + 1",
		},
		{
			name:  "multiline input",
			input: "3 + 1 / \n 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lox{}
			l.run(tt.input)

			// Currently run() just prints, so we're just verifying it doesn't panic
			assert.False(t, l.hadError, "expected hadError to remain false")
		})
	}
}
