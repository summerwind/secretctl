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

type GPGStorage struct {
	Config *config.GPGStorageConfig
}

func NewGPGStorage(c *config.GPGStorageConfig) (*GPGStorage, error) {
	return &GPGStorage{Config: c}, nil
}

func (s *GPGStorage) ReadSecret(p string) ([]byte, error) {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s\n", p)
	}

	if s.Config.Passphrase == "" {
		fmt.Print("Please enter the passphrase to unlock the secret key: ")
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return nil, err
		}
		fmt.Println("")

		s.Config.Passphrase = string(pass)
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	cmdOpts := []string{
		"--batch",
		"--no-tty",
		"--passphrase-fd", "0",
		"--pinentry-mode", "loopback",
		"-d", p,
	}

	cmd := exec.Command("gpg", cmdOpts...)
	cmd.Stdin = bytes.NewBuffer([]byte(s.Config.Passphrase))
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("gpg error: %s", stderr)
	}

	return stdout.Bytes(), nil
}

func (s *GPGStorage) WriteSecret(p string, data []byte) error {
	dir := filepath.Dir(p)

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("unable to create directory: %s\n", dir)
		}
	}

	if len(s.Config.Recipents) == 0 {
		return errors.New("no recipents of encrypted file")
	}

	cmdOpts := []string{"-e"}
	for _, r := range s.Config.Recipents {
		cmdOpts = append(cmdOpts, "-r", r)
	}

	file, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	stderr := bytes.NewBuffer([]byte{})

	cmd := exec.Command("gpg", cmdOpts...)
	cmd.Stdin = bytes.NewBuffer(data)
	cmd.Stdout = file
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("gpg error: %s", stderr)
	}

	return nil
}
