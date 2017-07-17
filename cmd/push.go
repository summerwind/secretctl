package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/summerwind/secretctl/config"
	"github.com/summerwind/secretctl/storage"
)

func NewPushCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "push",
		Short: "Push secret files and secret environment variables to remote storage",
		Long:  "Push secret files and secret environment variables to remote storage.",
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

	err = bindFlags(flags, c)
	if err != nil {
		return err
	}

	s, err := storage.NewStorage(c)
	if err != nil {
		return err
	}

	for path, secret := range c.Files {
		buf, err := ReadSecret(c.NormalizePath(path), false)
		if err != nil {
			return err
		}

		err = s.WriteSecret(secret, buf)
		if err != nil {
			if err == storage.ErrPullOnly || err == storage.ErrUnsupported {
				fmt.Printf("Skipped: %s (%s)\n", path, err)
				continue
			}

			return err
		}

		fmt.Printf("Pushed: %s\n", path)
	}

	for name, secret := range c.EnvVars {
		buf, err := ReadSecret(name, true)
		if err != nil {
			return err
		}

		err = s.WriteSecret(secret, buf)
		if err != nil {
			if err == storage.ErrPullOnly || err == storage.ErrUnsupported {
				fmt.Printf("Skipped: %s (%s)\n", name, err)
				continue
			}

			return err
		}

		fmt.Printf("Pushed: %s\n", name)
	}

	return nil
}
