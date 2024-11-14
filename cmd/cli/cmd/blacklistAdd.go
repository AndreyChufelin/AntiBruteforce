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

// blacklistAddCmd represents the blacklistAdd command
var blacklistAddCmd = &cobra.Command{
	Use:   "add",
	Short: "adds ip to blacklist",
	Long:  `adds ip to blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := iplist.BlacklistAdd(context.TODO(), &pbiplist.ListRequest{Ip: subnet})
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
