/*
Copyright © 2019 Netskope
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Constants used for invoke command and related flags
const (
	InvokeCmd      = "invoke"
	InvokeCmdShort = "Invoke a call against the service"
	InvokeCmdLong  = "Connect to the service and invoke an action"

	ClientHostFlagUsage = "the host to connect to"
	ClientPortFlagUsage = "the port to connect to"

	ClientHostFlag = "chost"
	ClientPortFlag = "cport"

	ClientHostDefaultValue = "0.0.0.0"
	ClientPortDefaultValue = 12345
)

var (
	// the command
	invokeCmd = &cobra.Command{
		Use:   InvokeCmd,
		Short: InvokeCmdShort,
		Long:  InvokeCmdLong,
		Args:  cobra.NoArgs,
	}
)

func init() {
	rootCmd.AddCommand(invokeCmd)

	// command flags
	invokeCmd.PersistentFlags().String(
		ClientHostFlag,
		ClientHostDefaultValue,
		ClientHostFlagUsage,
	)
	invokeCmd.PersistentFlags().Int32(
		ClientPortFlag,
		ClientPortDefaultValue,
		ClientPortFlagUsage,
	)

	err := viper.BindPFlags(invokeCmd.PersistentFlags())
	if err != nil {
		bail(err)
	}
}

func bail(err error) {
	fmt.Println(err)
	os.Exit(1)
}
