// +build !darwin

package storage

type KeychainStorage struct {
}

func NewKeychainStorage() (*KeychainStorage, error) {
	storage := KeychainStorage{}

	return &storage, nil
}

func (s *KeychainStorage) ReadSecret(label string) ([]byte, error) {
	return nil, ErrUnsupported
}

func (s *KeychainStorage) WriteSecret(label string, data []byte) error {
	return ErrUnsupported
}
