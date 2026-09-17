package source

import (
	"os"
	"strings"
	"testing"
)

func TestGetExtension(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{name: "monkey source file", file: "main.mx", expected: "mx"},
		{name: "source file with path", file: "/home/user/projects/monkey/main.mx", expected: "mx"},
		{name: "nested source file", file: "examples/hello/main.mx", expected: "mx"},
		{name: "go file", file: "main.go", expected: "go"},
		{name: "text file", file: "README.txt", expected: "txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extension := getExtension(tt.file)
			if extension != tt.expected {
				t.Errorf("getExtension(%q) = %q, expected %q", tt.file, extension, tt.expected)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		content  string
		expected string
	}{
		{name: "invalid extension", file: "main.go", content: `let x = 5;`, expected: `invalid file extension: ."go", expected .mx`},
		{name: "file does not exist", file: "missing.mx", expected: "no such file or directory"},
		{name: "parser error", file: "main.mx", content: `let = 5;`, expected: "parser errors:"},
		{name: "valid source file", file: "main.mx", content: `let x = 5;`, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.content != "" {
				err := os.WriteFile(tt.file, []byte(tt.content), 0644)
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tt.file)
			}
			originalArgs := os.Args
			defer func() {
				os.Args = originalArgs
			}()
			os.Args = []string{"monkey", tt.file}
			err := Run()
			if tt.expected == "" {
				if err != nil {
					t.Errorf("Run() returned unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("Run() expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Run() error = %q, expected to contain %q", err.Error(), tt.expected)
			}
		})
	}
}
