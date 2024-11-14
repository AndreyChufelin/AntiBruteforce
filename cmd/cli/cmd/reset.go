package cmd

import (
	"context"

	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var (
	login string
	ip    string
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "resets bucket",
	Long:  `resets bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := limiter.Clear(context.TODO(), &pbratelimter.ClearRequest{Login: login, Ip: ip})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				cmd.PrintErrln(e.Message())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.PersistentFlags().StringVar(&login, "login", "", "bucket login")
	resetCmd.PersistentFlags().StringVar(&ip, "ip", "", "bucket ip")
	resetCmd.MarkFlagsOneRequired("login", "ip")
}
