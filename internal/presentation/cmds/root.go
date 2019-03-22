package cmds

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Resource monitoring commands",
}

// RootCommand .
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(serveCommand())
	return rootCmd
}
