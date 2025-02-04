package build

import (
	"github.com/falcosecurity/dbg-go/pkg/build"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewBuildConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build dbg configs",
		RunE:  executeConfigs,
	}
	flags := cmd.Flags()
	flags.Bool("skip-existing", true, "whether to skip the build of drivers existing on S3")
	flags.Bool("publish", false, "whether artifacts must be published on S3")
	flags.Bool("ignore-errors", false, "whether to ignore build errors and go on looping on config files")
	flags.String("redirect-errors", "", "redirect build errors to the specified file")
	return cmd
}

func executeConfigs(_ *cobra.Command, _ []string) error {
	options := build.Options{
		Options:        root.LoadRootOptions(),
		SkipExisting:   viper.GetBool("skip-existing"),
		Publish:        viper.GetBool("publish"),
		IgnoreErrors:   viper.GetBool("ignore-errors"),
		RedirectErrors: viper.GetString("redirect-errors"),
	}
	return build.Run(options)
}
