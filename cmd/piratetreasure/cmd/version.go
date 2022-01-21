/*
Copyright © 2019 Netskope
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/netskope/piratetreasure/internal/build"
)

// Constants used for version command
const (
	VersionCmd      = "version"
	VersionCmdShort = "Show version"
)

var (
	//the command
	versionCmd = &cobra.Command{
		Use:   VersionCmd,
		Short: VersionCmdShort,
		RunE:  versionCmdFunc(),
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionCmdFunc() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s:\n Version: %s\n GitSha: %s\n Built: %s by %s\n",
			build.AppName, build.Version, build.GitSha, build.BuildTime, build.BuiltBy)

		return nil
	}
}
