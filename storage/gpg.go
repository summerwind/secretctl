package storage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/summerwind/secretctl/config"
)

// GPGStorage represents a local storage with GPG encription.
type GPGStorage struct {
	Config *config.GPGStorageConfig
}

// NewGPGStorage returns a GPGStorage with specified configuration.
func NewGPGStorage(c *config.GPGStorageConfig) (*GPGStorage, error) {
	return &GPGStorage{Config: c}, nil
}

// ReadSecret decrypts a encrypted file with the specified path by
// using gpg command.
func (s *GPGStorage) ReadSecret(p string) ([]byte, error) {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("File does not exist: %s", p)
	}

	if s.Config.Passphrase == "" {
		fmt.Print("Please enter the passphrase to unlock the secret key: ")
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return nil, fmt.Errorf("Unable to read password: %s", err)
		}
		fmt.Println("")

		s.Config.Passphrase = string(pass)
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	gpgCmd := s.Config.Command
	if gpgCmd == "" {
		gpgCmd = "gpg"
	}

	gpgOpts := []string{
		"--batch",
		"--no-tty",
		"--passphrase-fd", "0",
		"--pinentry-mode", "loopback",
		"-d", p,
	}

	cmd := exec.Command(gpgCmd, gpgOpts...)
	cmd.Stdin = bytes.NewBuffer([]byte(s.Config.Passphrase))
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Command error:\n%s", stderr)
	}

	return stdout.Bytes(), nil
}

// WriteSecret encrypts a secret to the specified path by using gpg
// command.
func (s *GPGStorage) WriteSecret(p string, data []byte) error {
	dir := filepath.Dir(p)

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("Unable to create directory: %s", dir)
		}
	}

	if len(s.Config.Recipents) == 0 {
		return errors.New("No recipents for encrypted file")
	}

	file, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Unable to open file: %s - %s", p, err)
	}
	defer file.Close()

	stderr := bytes.NewBuffer([]byte{})

	gpgCmd := s.Config.Command
	if gpgCmd == "" {
		gpgCmd = "gpg"
	}

	gpgOpts := []string{"-e"}
	for _, r := range s.Config.Recipents {
		gpgOpts = append(gpgOpts, "-r", r)
	}

	cmd := exec.Command(gpgCmd, gpgOpts...)
	cmd.Stdin = bytes.NewBuffer(data)
	cmd.Stdout = file
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Command error:\n%s", stderr)
	}

	return nil
}
