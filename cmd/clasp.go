package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vdemeester/clasp/config"
)

var (
	version   = "dev"
	gitsha    = "unknown"
	buildDate = ""
)

type tuckOptions struct {
	version bool
}

// NewRootCommand returns a new root command
func NewRootCommand() *cobra.Command {
	var opts tuckOptions

	cmd := &cobra.Command{
		Use:   "clasp CONFIG_FILE HOOK_DIR",
		Short: "a mini hook / rebuild configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("clasp requires 2 arguments exactly.")
			}
			return runClasp(args[0], args[1], opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.version, "version", false, "Print version and exit")
	return cmd
}

func runClasp(target, sourceDir string, opts tuckOptions) error {
	if opts.version {
		printVersion()
		return nil
	}
	configFile := config.NewFile(target, sourceDir)
	return configFile.Rebuild()
}

func printVersion() {
	fmt.Printf("clasp version %s (build: %s, date: %s)\n", version, gitsha, buildDate)
}
