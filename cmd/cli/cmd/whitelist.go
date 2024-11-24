package cmd

import (
	"github.com/spf13/cobra"
)

var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "avalible actions for whitelist [add delete]",
	Long:  `avalible actions for whitelist [add delete]`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.PrintErrln("must specify action [add delete]")
	},
}

func init() {
	rootCmd.AddCommand(whitelistCmd)
}
