package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
)

func TestWriteFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "runfromyaml_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		content  string
		filename string
		perm     os.FileMode
		wantErr  bool
	}{
		{
			name:     "write simple file",
			content:  "Hello, World!",
			filename: "test.txt",
			perm:     0644,
			wantErr:  false,
		},
		{
			name:     "write file with different permissions",
			content:  "#!/bin/bash\necho 'test'",
			filename: "script.sh",
			perm:     0755,
			wantErr:  false,
		},
		{
			name:     "write empty file",
			content:  "",
			filename: "empty.txt",
			perm:     0644,
			wantErr:  false,
		},
		{
			name:     "write file with special characters",
			content:  "Special chars: äöü ñ 中文 🚀",
			filename: "special.txt",
			perm:     0644,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tt.filename)

			// Test should not panic for valid inputs
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("WriteFile() panicked: %v", r)
				}
			}()

			WriteFile(tt.content, filePath, tt.perm)

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("File was not created: %s", filePath)
					return
				}

				// Verify file content
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read created file: %v", err)
					return
				}

				if string(content) != tt.content {
					t.Errorf("File content = %q, want %q", string(content), tt.content)
				}

				// Verify file permissions
				info, err := os.Stat(filePath)
				if err != nil {
					t.Errorf("Failed to get file info: %v", err)
					return
				}

				if info.Mode().Perm() != tt.perm {
					t.Errorf("File permissions = %o, want %o", info.Mode().Perm(), tt.perm)
				}
			}
		})
	}
}

func TestWriteFileWithEnvVar(t *testing.T) {
	// Set up environment variable
	tempDir, err := os.MkdirTemp("", "runfromyaml_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	os.Setenv("TEST_DIR", tempDir)
	defer os.Unsetenv("TEST_DIR")

	content := "Test content"
	filePath := "$TEST_DIR/envtest.txt"

	WriteFile(content, filePath, 0644)

	// Verify file was created in the expanded path
	expandedPath := filepath.Join(tempDir, "envtest.txt")
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		t.Errorf("File was not created at expanded path: %s", expandedPath)
	}

	// Verify content
	readContent, err := os.ReadFile(expandedPath)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("File content = %q, want %q", string(readContent), content)
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		index    int
		expected []string
	}{
		{
			name:     "remove first element",
			slice:    []string{"a", "b", "c", "d"},
			index:    0,
			expected: []string{"b", "c", "d"},
		},
		{
			name:     "remove middle element",
			slice:    []string{"a", "b", "c", "d"},
			index:    2,
			expected: []string{"a", "b", "d"},
		},
		{
			name:     "remove last element",
			slice:    []string{"a", "b", "c", "d"},
			index:    3,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "remove from single element slice",
			slice:    []string{"only"},
			index:    0,
			expected: []string{},
		},
		{
			name:     "remove from two element slice",
			slice:    []string{"first", "second"},
			index:    1,
			expected: []string{"first"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the slice to avoid modifying the test data
			testSlice := make([]string, len(tt.slice))
			copy(testSlice, tt.slice)

			result := Remove(testSlice, tt.index)

			if len(result) != len(tt.expected) {
				t.Errorf("Result length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Result[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestConfig(t *testing.T) {
	// Config() is currently empty, so we just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Config() panicked: %v", r)
		}
	}()

	Config()
}

// Test helper functions that might be used in the functions package
func TestPrintColor(t *testing.T) {
	// This is a basic test to ensure PrintColor doesn't panic
	// In a real scenario, you might want to capture output and verify it
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintColor() panicked: %v", r)
		}
	}()

	// Test with different color attributes
	PrintColor(color.FgRed, "test", "message")
	PrintColor(color.FgGreen, "success", "operation completed")
	PrintColor(color.FgYellow, "warning", "this is a warning")
}

// Benchmark tests
func BenchmarkRemove(b *testing.B) {
	slice := make([]string, 1000)
	for i := range slice {
		slice[i] = "item"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testSlice := make([]string, len(slice))
		copy(testSlice, slice)
		Remove(testSlice, 500)
	}
}

func BenchmarkWriteFile(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content := "This is test content for benchmarking file write operations."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("bench_test_%d.txt", i))
		WriteFile(content, filename, 0644)
	}
}
