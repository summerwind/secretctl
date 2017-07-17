package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/summerwind/secretctl/config"
	"github.com/summerwind/secretctl/storage"
)

func NewPullCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull secret files from remote storage",
		Long:  "Pull secret files from remote storage.",
		RunE:  runPullCommand,
	}

	return cmd
}

func runPullCommand(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	cp, err := flags.GetString("config")
	if err != nil {
		return err
	}

	c, err := config.Load(cp)
	if err != nil {
		return err
	}

	err = bindFlags(flags, c)
	if err != nil {
		return err
	}

	vault, err := storage.NewVaultStorage(c.Storage.Vault)
	if err != nil {
		return err
	}

	keychain, err := storage.NewKeychainStorage()
	if err != nil {
		return err
	}

	gpg, err := storage.NewGPGStorage(c.Storage.GPG)
	if err != nil {
		return err
	}

	for path, s := range c.Files {
		var (
			buf []byte
			err error
		)

		switch {
		case s.Vault != nil:
			buf, err = vault.ReadSecret(s.Vault.Path)
		case s.Keychain != nil:
			buf, err = keychain.ReadSecret(s.Keychain.Label)
		case s.GPG != nil:
			buf, err = gpg.ReadSecret(NormalizePath(cp, s.GPG.Path))
		default:
			err = errors.New("Storage configuration required")
		}

		_, err = WriteSecret(NormalizePath(cp, path), buf, false)
		if err != nil {
			if err == storage.Unsupported {
				fmt.Printf("[File] Skipped: %s (unsupported)\n", path)
				continue
			}
			return err
		}

		fmt.Printf("[File] Pulled: %s\n", path)
	}

	return nil
}
