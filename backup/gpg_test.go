package backup

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEncryptFileWithGPG(t *testing.T) {
	// Check if gpg is installed
	if _, err := exec.LookPath("gpg"); err != nil {
		t.Skip("gpg not found in PATH, skipping test")
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testgpg")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy file to encrypt
	inFile := filepath.Join(tmpDir, "test.txt")
	file, err := os.Create(inFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.WriteString("hello world")
	file.Close()

	// Encrypt the file
	outFile := filepath.Join(tmpDir, "test.txt.gpg")
	password := "testpassword"
	err = encryptFileWithGPG(inFile, outFile, password)
	if err != nil {
		t.Fatalf("encryptFileWithGPG failed: %v", err)
	}

	// Check if the encrypted file exists
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Errorf("Encrypted file was not created: %s", outFile)
	}

	//Decrypt the file to verify
	decryptedFile := filepath.Join(tmpDir, "decrypted.txt")
	cmd := exec.Command("gpg", "--batch", "--yes", "--passphrase", password, "-o", decryptedFile, "-d", outFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to decrypt file: %v", err)
	}

	// Check if the decrypted file exists
	if _, err := os.Stat(decryptedFile); os.IsNotExist(err) {
		t.Errorf("Decrypted file was not created: %s", decryptedFile)
	}

	// Verify the content of the decrypted file
	content, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}
	expectedContent := "hello world"
	if string(content) != expectedContent {
		t.Errorf("Decrypted file content mismatch. Got: %s, Expected: %s", string(content), expectedContent)
	}
}
