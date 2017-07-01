package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/summerwind/secretctl/cmd"
)

var (
	VERSION string = "2.0.0"
	COMMIT  string = "HEAD"
)

func main() {
	var cli = &cobra.Command{
		Use:   "secretctl <command>",
		Short: "Yet another secret management utility",
		Long:  "Yet another secret management utility.",
		RunE:  run,
	}

	flags := cli.PersistentFlags()
	flags.StringP("config", "c", ".secret.yml", "Path of configuration file")
	flags.Bool("version", false, "Display version information and exit")
	flags.Bool("help", false, "Display this help and exit")

	flags.String("vault-token", "", "The authentication token for the Vault server")
	flags.String("vault-addr", "", "The address of the Vault server")
	flags.String("vault-ca-cert", "", "Path to a PEM encoded CA cert file to use to verify the Vault server certificate")
	flags.String("vault-ca-path", "", "Path to a directory of PEM encoded CA cert files to verify the Vault server certificate")
	flags.String("vault-client-cert", "", "Path to a PEM encoded client certificate for TLS authentication to the Vault server")
	flags.String("vault-client-key", "", "Path to an unencrypted PEM encoded private key matching the client certificate")
	flags.Bool("vault-tls-skip-verify", false, "Do not verify TLS certificate")

	flags.StringSlice("gpg-recipent", []string{}, "Users who can decrypt files")
	flags.String("gpg-passphrase", "", "The passphrase of the GPG key")
	flags.String("gpg-command", "", "Path to a file to use as gpg command")

	cli.AddCommand(cmd.NewPushCommand())
	cli.AddCommand(cmd.NewPullCommand())
	cli.AddCommand(cmd.NewExecCommand())

	cli.SilenceUsage = true
	cli.SilenceErrors = true

	err := cli.Execute()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	flags := cmd.PersistentFlags()

	v, err := flags.GetBool("version")
	if err != nil {
		return err
	}

	if v {
		version()
		os.Exit(0)
	}

	return nil
}

func version() {
	fmt.Printf("Version: %s (%s)\n", VERSION, COMMIT)
}
