package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Config represents a configuration of secretctl.This includes
// configuration for storage, secret files and secret env vars.
type Config struct {
	Storage *StorageConfig     `json:"storage"`
	Files   map[string]*Secret `json:"files"`
	EnvVars map[string]*Secret `json:"env_vars"`
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
	PullOnly bool            `json:"pull_only"`
	Vault    *VaultOption    `json:"vault"`
	Keychain *KeychainOption `json:"keychain"`
	GPG      *GPGOption      `json:"gpg"`
}

// VaultOption represents a option for Vault storage.
type VaultOption struct {
	Path string `json:"path"`
}

// KeychainOption represents a option for Vault storage.
type KeychainOption struct {
	Label string `json:"label"`
}

// GPGOption represents a option for GPG storage.
type GPGOption struct {
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

	return &c, nil
}
