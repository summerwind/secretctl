package cmd

import (
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

	s, err := storage.NewStorage(c)
	if err != nil {
		return err
	}

	for path, secret := range c.Files {
		buf, err := s.ReadSecret(secret)
		if err != nil {
			if err == storage.ErrUnsupported {
				fmt.Printf("Skipped: %s (%s)\n", path, err)
				continue
			}

			return err
		}

		_, err = WriteSecret(c.NormalizePath(path), buf, false)
		if err != nil {
			return err
		}

		fmt.Printf("Pulled: %s\n", path)
	}

	return nil
}
