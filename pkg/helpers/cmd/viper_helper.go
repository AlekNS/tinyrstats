package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigureViper is need for preparing variables by an external configuration file.
func ConfigureViper(svcName, configName string, cmd *cobra.Command, viper *viper.Viper) *viper.Viper {
	var flagConfig = cmd.Flag("config-file")

	if len(flagConfig.Value.String()) == 0 {
		// if config not specified
		viper.SetConfigName(configName)
		viper.AddConfigPath("/etc/" + svcName)
		viper.AddConfigPath("$HOME/." + svcName)
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(flagConfig.Value.String())
	}

	var err = viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	return viper
}
