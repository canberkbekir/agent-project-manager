package cmd

import (
	"fmt"
	"os"

	"agent-project-manager/internal/agentctl/app"
	"agent-project-manager/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "agentctl",
	Short:        "CLI for agentd",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd == cmd.Root() {
			return nil
		}
		if h, _ := cmd.Flags().GetBool("help"); h {
			return nil
		}
		if cmd.Name() == "help" {
			return nil
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}

		a := &app.App{
			Cfg: cfg,
			Out: cmd.OutOrStdout(),
			Err: cmd.ErrOrStderr(),
		}
		cmd.SetContext(app.WithApp(cmd.Context(), a))
		return nil
	},
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
