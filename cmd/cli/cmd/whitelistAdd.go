package cmd

import (
	"context"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var whitelistAddCmd = &cobra.Command{
	Use:   "add",
	Short: "adds ip to whitelist",
	Long:  `adds ip to whitelist`,
	Run: func(cmd *cobra.Command, _ []string) {
		_, err := iplist.WhitelistAdd(context.TODO(), &pbiplist.ListRequest{Ip: subnet})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				cmd.PrintErrln(e.Message())
			}
		}
	},
}

func init() {
	whitelistCmd.AddCommand(whitelistAddCmd)

	whitelistAddCmd.Flags().StringVar(&subnet, "subnet", "", "subnet to add to whitelist")
	whitelistAddCmd.MarkFlagRequired("subnet")
}
