// +build darwin

package storage

import (
	"errors"
	"fmt"

	keychain "github.com/keybase/go-keychain"
)

const (
	SERVICE_NAME = "secretctl"
)

type KeychainStorage struct {
}

func NewKeychainStorage() (*KeychainStorage, error) {
	storage := KeychainStorage{}

	return &storage, nil
}

func (s *KeychainStorage) ReadSecret(label string) ([]byte, error) {
	if label == "" {
		return nil, errors.New("label is required")
	}

	query := s.newItem(label)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return nil, err
	}

	return results[0].Data, nil
}

func (s *KeychainStorage) WriteSecret(label string, data []byte) error {
	if label == "" {
		return errors.New("label is required")
	}

	query := s.newItem(label)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnAttributes(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return err
	}

	item := s.newItem(label)
	item.SetData(data)
	item.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)

	if len(results) == 0 {
		err := keychain.AddItem(item)
		if err != nil {
			return err
		}
	} else {
		err := keychain.UpdateItem(query, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *KeychainStorage) newItem(label string) keychain.Item {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(SERVICE_NAME)
	item.SetAccount(label)
	item.SetLabel(fmt.Sprintf("%s: %s", SERVICE_NAME, label))

	return item
}
