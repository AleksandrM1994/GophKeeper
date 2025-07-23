package cobra_cli

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "GophKeeper",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra_cli.Command, args []string) { },
	}
	return rootCmd
}
