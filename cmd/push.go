package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/summerwind/secretctl/config"
	"github.com/summerwind/secretctl/storage"
)

func NewPushCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "push",
		Short: "Update remote secrets based on local secrets",
		Long:  "Update remote secrets based on local secrets.",
		RunE:  runPushCommand,
	}

	return cmd
}

func runPushCommand(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	cp, err := flags.GetString("config")
	if err != nil {
		return err
	}

	c, err := config.Load(cp)
	if err != nil {
		return err
	}

	vault, err := storage.NewVaultStorage(c.Storage.Vault)
	if err != nil {
		return err
	}

	gpg, err := storage.NewGPGStorage(c.Storage.GPG)
	if err != nil {
		return err
	}

	for path, s := range c.Files {
		if s.PullOnly {
			fmt.Printf("[File] Skipped: %s (pull only)\n", path)
			continue
		}

		buf, err := ReadSecret(NormalizePath(cp, path), false)
		if err != nil {
			return err
		}

		switch {
		case s.Vault != nil:
			err = vault.WriteSecret(s.Vault.Path, buf)
		case s.GPG != nil:
			err = gpg.WriteSecret(NormalizePath(cp, s.GPG.Path), buf)
		default:
			err = errors.New("Storage configuration is required")
		}

		fmt.Printf("[File] Pushed: %s\n", path)
	}

	for name, s := range c.EnvVars {
		if s.PullOnly {
			fmt.Printf("[EnvVar] Skipped: %s (pull only)\n", name)
			continue
		}

		buf, err := ReadSecret(name, true)
		if err != nil {
			return err
		}

		switch {
		case s.Vault != nil:
			err = vault.WriteSecret(s.Vault.Path, buf)
		case s.GPG != nil:
			err = gpg.WriteSecret(NormalizePath(cp, s.GPG.Path), buf)
		default:
			err = errors.New("Storage configuration is required")
		}

		fmt.Printf("[EnvVar] Pushed: %s\n", name)
	}

	return nil
}
