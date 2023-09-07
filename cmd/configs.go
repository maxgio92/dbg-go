package cmd

import (
	"github.com/fededp/dbg-go/cmd/build"
	"github.com/fededp/dbg-go/cmd/cleanup"
	"github.com/fededp/dbg-go/cmd/generate"
	"github.com/fededp/dbg-go/cmd/stats"
	"github.com/fededp/dbg-go/cmd/validate"
	"github.com/spf13/cobra"
)

var (
	configsCmd = &cobra.Command{
		Use:   "configs",
		Short: "Work with local dbg configs",
	}
)

func init() {
	// Subcommands
	configsCmd.AddCommand(generate.NewGenerateConfigsCmd())
	configsCmd.AddCommand(cleanup.NewCleanupConfigsCmd())
	configsCmd.AddCommand(validate.NewValidateConfigsCmd())
	configsCmd.AddCommand(stats.NewStatsConfigsCmd())
	configsCmd.AddCommand(build.NewBuildConfigsCmd())
}
