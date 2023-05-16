/*
Copyright © 2023 Giancarlo Espinoza <gespiza@gmail.com>

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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goblin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.PersistentFlags().StringP("asset", "a", "BTC/USD", "asset to query; e.g. BTC/USD, ETH/USD")
	rootCmd.PersistentFlags().StringP("exchange", "x", "BINANCE", "exchange to pull data from; e.g. BINANCE, COINBASE")
	rootCmd.PersistentFlags().StringP("indicator", "i", "supertrend", "name of the indicator to use")
	rootCmd.PersistentFlags().StringP("interval", "t", "1h", "time interval; e.g. 5m 10m 1h 1d 1week")
	rootCmd.PersistentFlags().StringP("log-level", "v", "DEBUG", "log level: INFO, WARN, ERROR, DEBUG")
	rootCmd.PersistentFlags().StringP("start-date", "s", "", "start date to query a chart from; e.g 2022-08-06")
	rootCmd.PersistentFlags().StringP("end-date", "e", "", "end date to query a chart from; e.g. 2022-08-31")
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()

	// rootCmd.MarkFlagRequired("asset")
	// rootCmd.MarkFlagRequired("indicator")
}
