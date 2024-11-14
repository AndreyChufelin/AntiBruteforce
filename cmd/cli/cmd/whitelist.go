/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// whitelistCmd represents the whitelist command
var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "avalible actions for whitelist [add delete]",
	Long:  `avalible actions for whitelist [add delete]`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln("must specify action [add delete]")
	},
}

func init() {
	rootCmd.AddCommand(whitelistCmd)
}
