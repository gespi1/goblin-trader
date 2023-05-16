package cmd

import (
	"goblin-trader/pkg/apis/twelvedata"
	"goblin-trader/pkg/common"
	"goblin-trader/pkg/indicators"
	"goblin-trader/pkg/plotty"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chartCmd represents the chart command
var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Querys and charts historical data using an indicator; plots the data and renders a csv",
	Long:  `CLI tool to get quick historical data based on a indicator(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.GetViper()
		common.SetLogLevel(v.GetString("log-level"))
		twelve := twelvedata.Init(v)
		series, err := twelve.TimeSeries()
		if err != nil {
			log.Fatalf("not able to get dataframe from timeseries: %v", err)
		}
		// traditional SuperTrend lookback and multiplier values are 10 and 3 respectively
		ub, lb, superTrend := indicators.SuperTrend(series, 10, 3)
		plotty.SuperTrend(series, ub, lb, superTrend)
	},
}

func init() {
	rootCmd.AddCommand(chartCmd)
}
