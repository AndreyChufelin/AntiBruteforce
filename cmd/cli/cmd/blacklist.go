package cmd

import (
	"github.com/spf13/cobra"
)

var blacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "avalible actions for blacklist [add delete]",
	Long:  `avalible actions for blacklist [add delete]`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.PrintErrln("must specify action [add delete]")
	},
}

func init() {
	rootCmd.AddCommand(blacklistCmd)
}
