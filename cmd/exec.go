package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/summerwind/secretctl/config"
	"github.com/summerwind/secretctl/storage"
)

func NewExecCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "exec",
		Short: "Pull secret environment variables from remote storage and run specified command",
		Long:  "Pull secret environment variables from remote storage and run specified command.",
		RunE:  runExecCommand,
	}

	return cmd
}

func runExecCommand(cmd *cobra.Command, args []string) error {
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

	for name, secret := range c.EnvVars {
		buf, err := s.ReadSecret(secret)
		if err != nil {
			if err == storage.ErrUnsupported {
				fmt.Printf("Skipped: %s (%s)\n", name, err)
				continue
			}

			return err
		}

		_, err = WriteSecret(name, buf, true)
		if err != nil {
			return err
		}
	}

	ecmd := exec.Command(args[0], args[1:]...)
	ecmd.Stdin = os.Stdin
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	ecmd.Run()

	return nil
}
