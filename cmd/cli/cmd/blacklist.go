/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// blacklistCmd represents the blacklist command
var blacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "avalible actions for blacklist [add delete]",
	Long:  `avalible actions for blacklist [add delete]`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln("must specify action [add delete]")
	},
}

func init() {
	rootCmd.AddCommand(blacklistCmd)
}
