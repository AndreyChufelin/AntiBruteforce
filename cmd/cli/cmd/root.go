package cmd

import (
	"fmt"
	"log"
	"os"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cfgFile string
	limiter pbratelimter.RatelimiterClient
	iplist  pbiplist.IPListServiceClient
	subnet  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "antibruteforce-cli",
	Short: "antibruteforce-cli tool for administrating anti-bruteforce service",
	Long:  `antibruteforce-cli tool for administrating anti-bruteforce service`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed initialize client", err)
	}

	limiter = pbratelimter.NewRatelimiterClient(conn)
	iplist = pbiplist.NewIPListServiceClient(conn)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
