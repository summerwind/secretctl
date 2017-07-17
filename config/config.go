package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
)

// Config represents a configuration of secretctl.This includes
// configuration for storage, secret files and secret env vars.
type Config struct {
	BasePath string             `json:"base_path"`
	Storage  *StorageConfig     `json:"storage"`
	Files    map[string]*Secret `json:"files"`
	EnvVars  map[string]*Secret `json:"env_vars"`
}

// StorageConfig represents a configuration for storage. This includes
// configuration for Vault and GPG.
type StorageConfig struct {
	Vault *VaultStorageConfig `json:"vault"`
	GPG   *GPGStorageConfig   `json:"gpg"`
}

// VaultStorageConfig represents a configuration for Vault.
type VaultStorageConfig struct {
	Token         string `json:"token"`
	Addr          string `json:"addr"`
	CACert        string `json:"ca_cert"`
	CAPath        string `json:"ca_path"`
	ClientCert    string `json:"client_cert"`
	ClientKey     string `json:"client_key"`
	TLSSkipVerify bool   `json:"tls_skip_verify"`
}

// GPGStorageConfig represents configuration for GPG.
type GPGStorageConfig struct {
	Recipents  []string `json:"recipents"`
	Passphrase string   `json:"passphrase"`
	Command    string   `json:"command"`
}

// Secret represents secret file or secret env var.
type Secret struct {
	PullOnly bool           `json:"pull_only"`
	Vault    *VaultParam    `json:"vault"`
	Keychain *KeychainParam `json:"keychain"`
	GPG      *GPGParam      `json:"gpg"`
}

// VaultParam represents a set of parameters for Vault storage.
type VaultParam struct {
	Path string `json:"path"`
}

// KeychainParam represents a set of parameters for Vault storage.
type KeychainParam struct {
	Label string `json:"label"`
}

// GPGParam represents a set of parameters for GPG storage.
type GPGParam struct {
	Path string `json:"path"`
}

// Load loads configuration file with specified path and returns
// a Config.
func Load(p string) (*Config, error) {
	var c Config

	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("Unable to load configuration file: %s", err)
	}

	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse configuration file: %s", err)
	}

	if c.BasePath == "" {
		c.BasePath = filepath.Dir(p)
	}

	return &c, nil
}

// NormalizePath returns a path string which is combined the given
// path and base path based on the current configuration.
func (c *Config) NormalizePath(p string) string {
	if p != "" && !filepath.IsAbs(p) {
		p = filepath.Join(c.BasePath, p)
	}

	return p
}
