package storage

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/summerwind/secretctl/config"
)

type VaultStorage struct {
	Config *config.VaultStorageConfig
	Client *vault.Client
}

func NewVaultStorage(c *config.VaultStorageConfig) (*VaultStorage, error) {
	vc := vault.DefaultConfig()
	vc.Address = c.Addr

	vtc := &vault.TLSConfig{
		CACert:     c.CACert,
		CAPath:     c.CAPath,
		ClientCert: c.ClientCert,
		ClientKey:  c.ClientKey,
		Insecure:   c.TLSSkipVerify,
	}

	err := vc.ConfigureTLS(vtc)
	if err != nil {
		return nil, err
	}

	err = vc.ReadEnvironment()
	if err != nil {
		return nil, err
	}

	client, err := vault.NewClient(vc)
	if err != nil {
		return nil, err
	}

	token := os.Getenv(vault.EnvVaultToken)
	if token == "" {
		client.SetToken(c.Token)
	}

	storage := VaultStorage{
		Config: c,
		Client: client,
	}

	return &storage, nil
}

func (s *VaultStorage) ReadSecret(p string) ([]byte, error) {
	if p[0] == '/' {
		p = p[1:]
	}

	secret, err := s.Client.Logical().Read(p)
	if err != nil {
		return nil, fmt.Errorf("Unable to read secret from Vault: %s", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("Secret does not exist at %s", p)
	}

	raw, ok := secret.Data["value"]
	if !ok {
		return nil, fmt.Errorf("No secret found at %s", p)
	}

	val, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("No secret found at %s", p)
	}

	return []byte(val), nil
}

func (s *VaultStorage) WriteSecret(p string, data []byte) error {
	if p[0] == '/' {
		p = p[1:]
	}

	vd := map[string]interface{}{
		"value": string(data),
	}

	_, err := s.Client.Logical().Write(p, vd)
	if err != nil {
		return nil, fmt.Errorf("Unable to write secret to Vault: %s", err)
	}

	return nil
}
