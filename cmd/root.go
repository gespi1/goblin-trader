/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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
	v := viper.GetViper()
	setLogLevel(v.GetString("log-level"))
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goblin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.PersistentFlags().StringP("asset", "a", "BTC/USD", "asset to query; e.g. BTC/USD, ETH/USD")
	rootCmd.PersistentFlags().StringP("exchange", "x", "BINANCE", "exchange to pull data from; e.g. BINANCE, COINBASE")
	rootCmd.PersistentFlags().StringP("indicator", "i", "supertrend", "name of the indicator to use")
	rootCmd.PersistentFlags().StringP("interval", "t", "1h", "time interval; e.g. 5m 10m 1h 1d 1w")
	rootCmd.PersistentFlags().StringP("log-level", "v", "DEBUG", "log level: INFO, WARN, ERROR, DEBUG")
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()

	// rootCmd.MarkFlagRequired("asset")
	// rootCmd.MarkFlagRequired("indicator")
}

func setLogLevel(loglevel string) {
	l := strings.ToLower(loglevel)
	if l == "info" {
		log.SetLevel(log.InfoLevel)
	} else if l == "warn" {
		log.SetLevel(log.WarnLevel)
	} else if l == "error" {
		log.SetLevel(log.ErrorLevel)
	} else if l == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if l == "fatal" {
		log.SetLevel(log.FatalLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
		log.Warn("no log level matched setting to default, ERROR")
	}
}
