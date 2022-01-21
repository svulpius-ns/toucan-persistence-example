/*
Copyright © 2019-2020 Netskope
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	config "github.com/netskope/go-kestrel/pkg/config/cfgmanage"

	"github.com/netskope/piratetreasure/internal/build"
)

// Constants used for commands and related flags
const (
	RootCmdLong  = "its cursed"
	RootCmdShort = "piratetreasure"
)

var (
	cfgFile string

	gViper = viper.New()

	rootCmd = &cobra.Command{
		Use:     build.AppName,
		Short:   RootCmdShort,
		Long:    RootCmdLong,
		Version: fmt.Sprintf("%s (%s)", build.Version, build.GitSha),
	}
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// add a --config global flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		fmt.Sprintf("config file (default is $HOME/.%s.yaml)", build.AppName))

	config.InitConfig(
		config.CobraCommand(rootCmd),
		config.AppConfiguration(gViper),
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Config(config.GlobalConfig).SetDefault(
		"appname", build.AppName,
	)

	if cfgFile != "" { // enable ability to specify config file via flag
		gViper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// add home directory as first search path for a config file
		gViper.AddConfigPath(home)
		// name of config file (without extension)
		gViper.SetConfigName(fmt.Sprintf(".%s", build.AppName))
	}

	// read in environment variables that match

	// use environment variables with the uppercase form of the appname as a prefix
	gViper.SetEnvPrefix(build.AppName)
	// for config key names which have a '-', match against '_' when searching the ENV
	gViper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// read in environment variables that match
	gViper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := gViper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", gViper.ConfigFileUsed())
	}

	config.ApplyConfig()
}
