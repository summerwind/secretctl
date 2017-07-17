package storage

import (
	"errors"

	"github.com/summerwind/secretctl/config"
)

var (
	ErrUnsupported = errors.New("Unsupported")
	ErrNoParameter = errors.New("No storage parameter found")
	ErrPullOnly    = errors.New("Pull only")
)

type Storage struct {
	Config   *config.Config
	Vault    *VaultStorage
	Keychain *KeychainStorage
	GPG      *GPGStorage
}

func NewStorage(c *config.Config) (*Storage, error) {
	vault, err := NewVaultStorage(c.Storage.Vault)
	if err != nil {
		return nil, err
	}

	keychain, err := NewKeychainStorage()
	if err != nil {
		return nil, err
	}

	gpg, err := NewGPGStorage(c.Storage.GPG)
	if err != nil {
		return nil, err
	}

	s := Storage{
		Config:   c,
		Vault:    vault,
		Keychain: keychain,
		GPG:      gpg,
	}

	return &s, nil
}

func (s *Storage) ReadSecret(secret *config.Secret) ([]byte, error) {
	var (
		buf []byte
		err error
	)

	switch {
	case secret.Vault != nil:
		buf, err = s.Vault.ReadSecret(secret.Vault.Path)
	case s.Keychain != nil:
		buf, err = s.Keychain.ReadSecret(secret.Keychain.Label)
	case s.GPG != nil:
		buf, err = s.GPG.ReadSecret(s.Config.NormalizePath(secret.GPG.Path))
	default:
		err = ErrNoParameter
	}

	return buf, err
}

func (s *Storage) WriteSecret(secret *config.Secret, buf []byte) error {
	var err error

	if secret.PullOnly {
		return ErrPullOnly
	}

	switch {
	case secret.Vault != nil:
		err = s.Vault.WriteSecret(secret.Vault.Path, buf)
	case secret.Keychain != nil:
		err = s.Keychain.WriteSecret(secret.Keychain.Label, buf)
	case secret.GPG != nil:
		err = s.GPG.WriteSecret(s.Config.NormalizePath(secret.GPG.Path), buf)
	default:
		err = ErrNoParameter
	}

	return err
}
