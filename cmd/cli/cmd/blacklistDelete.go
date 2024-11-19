package cmd

import (
	"context"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var blacklistDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes ip from blacklist",
	Long:  `deletes ip from blacklist`,
	Run: func(cmd *cobra.Command, _ []string) {
		_, err := iplist.BlacklistDelete(context.TODO(), &pbiplist.ListRequest{Subnet: subnet})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				cmd.PrintErrln(e.Message())
			}
		}
	},
}

func init() {
	blacklistCmd.AddCommand(blacklistDeleteCmd)

	blacklistDeleteCmd.Flags().StringVar(&subnet, "subnet", "", "subnet to delete from blacklist")
	blacklistDeleteCmd.MarkFlagRequired("subnet")
}
