package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompressDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some dummy files in the temp directory
	file1, err := os.Create(filepath.Join(tmpDir, "file1.txt"))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file1.WriteString("hello")
	file1.Close()

	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	file2, err := os.Create(filepath.Join(subDir, "file2.txt"))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file2.WriteString("world")
	file2.Close()

	// Compress the directory
	destFile := filepath.Join(tmpDir, "test.tar.gz")
	err = compressDirectory(tmpDir, destFile)
	if err != nil {
		t.Fatalf("compressDirectory failed: %v", err)
	}

	// Check if the compressed file exists
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Errorf("Compressed file was not created: %s", destFile)
	}
}
