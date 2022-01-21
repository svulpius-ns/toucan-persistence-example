/*
Copyright © 2019-2022 Netskope
*/

package cmd

import (
	"context"
	"fmt"
	"os"

	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	empty "google.golang.org/protobuf/types/known/emptypb"

	apis "github.com/netskope-qe/toucan-base/api/proto/toucanbase"
)

// Constants used for health command and related flags
const (
	HealthCmd      = "health"
	HealthCmdShort = "Get health for the service"
	HealthCmdLong  = "Connect to the server and get the service health"
)

var (
	// the command
	healthCmd = &cobra.Command{
		Use:   HealthCmd,
		Short: HealthCmdShort,
		Long:  HealthCmdLong,
		Run:   healthCmdFunc(),
		Args:  cobra.NoArgs,
	}
)

func init() {
	invokeCmd.AddCommand(healthCmd)
}

func healthCmdFunc() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		host := viper.GetString(ClientHostFlag)
		port := viper.GetInt32(ClientPortFlag)

		if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port),
			grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
			bail(err)
		} else {
			cli := apis.NewHealthServiceClient(conn)
			health, err := cli.Health(context.Background(), &empty.Empty{})
			if err != nil {
				bail(err)
			}

			healthJSON, err := prettyjson.Marshal(health)
			if err != nil {
				bail(err)
			}

			fmt.Printf("%s %s\n", text.FgHiMagenta.Sprintf("Reply:"), healthJSON)
			os.Exit(0)
		}
	}
}
