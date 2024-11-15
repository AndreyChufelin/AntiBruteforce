package cmd

import (
	"context"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var blacklistAddCmd = &cobra.Command{
	Use:   "add",
	Short: "adds ip to blacklist",
	Long:  `adds ip to blacklist`,
	Run: func(cmd *cobra.Command, _ []string) {
		_, err := iplist.BlacklistAdd(context.TODO(), &pbiplist.ListRequest{Subnet: subnet})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				cmd.PrintErrln(e.Message())
			}
		}
	},
}

func init() {
	blacklistCmd.AddCommand(blacklistAddCmd)

	blacklistAddCmd.Flags().StringVar(&subnet, "subnet", "", "subnet to add to blacklist")
	blacklistAddCmd.MarkFlagRequired("subnet")
}
