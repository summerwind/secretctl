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
