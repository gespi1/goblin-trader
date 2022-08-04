/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goblin-trader",
	Short: "trading bot created in GOlang",
	Long: `CLI tool to get quick historical data based on a indicator(s) or
a Service to make trades on your behalf based on indicators.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.logrus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	runCmd.PersistentFlags().StringP("asset", "a", "BTC", "asset to assess; e.g. BTC, ETH. Additionally add an exhange to query; e.g. BTC/COINBASE, ETH/KRAKEN")
	runCmd.PersistentFlags().StringP("indicator", "i", "supertrend", "name of the indicator to use")
	runCmd.PersistentFlags().StringP("interval", "t", "1h", "time interval; e.g. 5m 10m 1h 1d 1w")
	viper.BindPFlags(rootCmd.PersistentFlags())
	// rootCmd.MarkFlagRequired("asset")
	// rootCmd.MarkFlagRequired("indicator")
}
