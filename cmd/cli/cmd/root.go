package cmd

import (
	"log"
	"net"
	"os"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cfgFile string
	limiter pbratelimter.RatelimiterClient
	iplist  pbiplist.IPListServiceClient
	subnet  string
)

var rootCmd = &cobra.Command{
	Use:   "antibruteforce-cli",
	Short: "antibruteforce-cli tool for administrating anti-bruteforce service",
	Long:  `antibruteforce-cli tool for administrating anti-bruteforce service`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	host, _ := os.LookupEnv("HOST")
	port, _ := os.LookupEnv("PORT")
	addr := net.JoinHostPort(host, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed initialize client", err)
	}

	limiter = pbratelimter.NewRatelimiterClient(conn)
	iplist = pbiplist.NewIPListServiceClient(conn)
}
