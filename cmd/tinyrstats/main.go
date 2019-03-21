package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alekns/tinyrstats/internal/presentation/cmds"
)

// ResourceStatsVersion .
const ResourceStatsVersion = "0.0.1"

var mainCommand = &cobra.Command{
	Use:     "tinyrstats",
	Short:   "Tiny Resource Statistics is a service for resource monitoring and gathering of some statistics",
	Version: ResourceStatsVersion}

func main() {
	viper.SetEnvPrefix("trs")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	mainFlags := mainCommand.PersistentFlags()
	mainFlags.String("config-file", "", "Config file")

	mainFlags.String("logging-level", "", "Default console level")
	viper.BindPFlag("logging.console.level", mainFlags.Lookup("logging-level"))

	mainCommand.AddCommand(cmds.RootCommand())

	if mainCommand.Execute() != nil {
		os.Exit(1)
	}
}
