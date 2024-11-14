/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// blacklistDeleteCmd represents the blacklistDelete command
var blacklistDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes ip from blacklist",
	Long:  `deletes ip from blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := iplist.BlacklistDelete(context.TODO(), &pbiplist.ListRequest{Ip: subnet})
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
