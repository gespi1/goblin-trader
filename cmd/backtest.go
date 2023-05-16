package cmd

import (
	"fmt"
	"goblin-trader/pkg/apis/twelvedata"
	"goblin-trader/pkg/common"
	"goblin-trader/pkg/indicators"
	"goblin-trader/pkg/strategies"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backtestCmd represents the backtest command
var backtestCmd = &cobra.Command{
	Use:   "backtest",
	Short: "Backtests trading rules and strategies",
	Long: `Backtests trading rules, indicator parameters and strategies against an indicator based on historical data. 
Returns a CSV, a plot and the results of the backtest`,
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.GetViper()
		common.SetLogLevel(v.GetString("log-level"))
		twelve := twelvedata.Init(v)
		series, err := twelve.TimeSeries()
		if err != nil {
			log.Fatalf("not able to get dataframe from timeseries: %v", err)
		}
		// traditional SuperTrend lookback and multiplier values are 10 and 3 respectively
		_, _, superTrend := indicators.SuperTrend(series, 10, 3)

		stopLoss := v.GetFloat64("stop-loss")
		takeProfit := v.GetFloat64("take-profit")
		startingCaptial := v.GetFloat64("starting-captial")

		prices := common.GetPricesSlice(series)
		trades := strategies.SupertrendStrategy(prices, superTrend)
		fmt.Println("TRADES")
		fmt.Println(trades)
		rules := strategies.TradingRules{
			StartingCapital: startingCaptial,
			StopLossPct:     stopLoss,
			TakeProfitPct:   takeProfit,
		}

		executedTrades, profit, capital := rules.SuperTrendExecuteTrades(prices, trades)
		fmt.Println("EXECTRADES")
		fmt.Println(executedTrades)
		fmt.Println("PROFIT")
		fmt.Println(profit)
		fmt.Println("CAPITAL")
		fmt.Println(capital)
	},
}

func init() {
	rootCmd.AddCommand(backtestCmd)

	backtestCmd.Flags().Float64("stop-loss", 2, "in percentage the price amount willing risk before triggering to sell at a loss (default 2%)")
	backtestCmd.Flags().Float64("take-profit", 5, "in percentage the price amount to take profits at (default 5%)")
	backtestCmd.Flags().Float64("starting-capital", 1000, "initial captial to start trade with")
	viper.BindPFlags(backtestCmd.Flags())
}
