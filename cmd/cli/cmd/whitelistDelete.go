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

// whitelistDeleteCmd represents the whitelistDelete command
var whitelistDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes ip from whitelist",
	Long:  `deletes ip from whitelist`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := iplist.WhitelistDelete(context.TODO(), &pbiplist.ListRequest{Ip: subnet})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				cmd.PrintErrln(e.Message())
			}
		}
	},
}

func init() {
	whitelistCmd.AddCommand(whitelistDeleteCmd)

	whitelistDeleteCmd.Flags().StringVar(&subnet, "subnet", "", "subnet to delete from whitelist")
	whitelistDeleteCmd.MarkFlagRequired("subnet")
}
