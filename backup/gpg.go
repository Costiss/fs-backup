package backup

import (
	"os"
	"os/exec"
)

func encryptFileWithGPG(inPath, outPath, password string) error {
	cmd := exec.Command("gpg", "--batch", "--yes", "--passphrase", password, "--symmetric", "--cipher-algo", "AES256", "-o", outPath, inPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
