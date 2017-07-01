package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/summerwind/secretctl/config"
)

func bindFlags(flags *pflag.FlagSet, c *config.Config) error {
	vsc := c.Storage.Vault
	gsc := c.Storage.GPG

	vaultToken, err := flags.GetString("vault-token")
	if err != nil {
		return err
	}

	if vaultToken != "" {
		vsc.Token = vaultToken
	}

	vaultAddr, err := flags.GetString("vault-addr")
	if err != nil {
		return err
	}

	if vaultAddr != "" {
		vsc.Addr = vaultAddr
	}

	vaultCACert, err := flags.GetString("vault-ca-cert")
	if err != nil {
		return err
	}

	if vaultCACert != "" {
		vsc.CACert = vaultCACert
	}

	vaultCAPath, err := flags.GetString("vault-ca-path")
	if err != nil {
		return err
	}

	if vaultCAPath != "" {
		vsc.CAPath = vaultCAPath
	}

	vaultClientCert, err := flags.GetString("vault-client-cert")
	if err != nil {
		return err
	}

	if vaultClientCert != "" {
		vsc.ClientCert = vaultClientCert
	}

	vaultClientKey, err := flags.GetString("vault-client-key")
	if err != nil {
		return err
	}

	if vaultClientKey != "" {
		vsc.ClientKey = vaultClientKey
	}

	vaultTLSSkipVerify, err := flags.GetBool("vault-tls-skip-verify")
	if err != nil {
		return err
	}

	if vaultTLSSkipVerify {
		vsc.TLSSkipVerify = vaultTLSSkipVerify
	}

	gpgRecipents, err := flags.GetStringSlice("gpg-recipent")
	if err != nil {
		return err
	}

	if len(gpgRecipents) > 0 {
		gsc.Recipents = gpgRecipents
	}

	gpgPassphrase, err := flags.GetString("gpg-passphrase")
	if err != nil {
		return err
	}

	if gpgPassphrase != "" {
		gsc.Passphrase = gpgPassphrase
	}

	gpgCommand, err := flags.GetString("gpg-command")
	if err != nil {
		return err
	}

	if gpgCommand != "" {
		gsc.Command = gpgCommand
	}

	fmt.Println(vsc)
	return nil
}

func NormalizePath(cp, fp string) string {
	dir := filepath.Dir(cp)

	if fp != "" && !filepath.IsAbs(fp) {
		fp = filepath.Join(dir, fp)
	}

	return fp
}

func ReadSecret(key string, env bool) ([]byte, error) {
	var (
		buf []byte
		err error
	)

	if env {
		ev := os.Getenv(key)
		if ev == "" {
			return nil, fmt.Errorf("environment variable does not exist: %s", key)
		}

		buf = []byte(ev)
	} else {
		buf, err = ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func WriteSecret(key string, data []byte, env bool) (int, error) {
	if env {
		err := os.Setenv(key, string(data))
		if err != nil {
			return 0, err
		}
	} else {
		dir := filepath.Dir(key)

		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				return 0, fmt.Errorf("unable to create directory: %s\n", dir)
			}
		}

		err = ioutil.WriteFile(key, data, 0600)
		if err != nil {
			return 0, err
		}
	}

	return len(data), nil
}
