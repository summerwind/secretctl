package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type Config struct {
	Storage *StorageConfig     `json:"storage"`
	Files   map[string]*Secret `json:"files"`
	EnvVars map[string]*Secret `json:"env_vars"`
}

type StorageConfig struct {
	Vault *VaultStorageConfig `json:"vault"`
	GPG   *GPGStorageConfig   `json:"gpg"`
}

type VaultStorageConfig struct {
	Token         string `json:"token"`
	Addr          string `json:"addr"`
	CACert        string `json:"ca_cert"`
	CAPath        string `json:"ca_path"`
	ClientCert    string `json:"client_cert"`
	ClientKey     string `json:"client_key"`
	TLSSkipVerify bool   `json:"tls_skip_verify"`
}

type GPGStorageConfig struct {
	Recipents  []string `json:"recipents"`
	Passphrase string   `json:"passphrase"`
	Command    string   `json:"command"`
}

type Secret struct {
	PullOnly bool         `json:"pull_only"`
	Vault    *VaultOption `json:"vault"`
	GPG      *GPGOption   `json:"gpg"`
}

type VaultOption struct {
	Path string `json:"path"`
}

type GPGOption struct {
	Path string `json:"path"`
}

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
