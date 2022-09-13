/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"goblin-trader/pkg/apis/twelvedata"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs to query historical data using an indicator",
	Long:  `CLI tool to get quick historical data based on a indicator(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.GetViper()
		twelve := twelvedata.Init(v)
		twelve.TimeSeries()
		// if err != nil {
		// 	log.Fatalf("not able to get dataframe from timeseries: %v", err)
		// }
		// fmt.Println(df)
		// twelvedata.SuperTrend(df)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	runCmd.PersistentFlags().StringP("start-date", "s", "", "start date to query a chart from; e.g 2022-08-06")
	runCmd.PersistentFlags().StringP("end-date", "e", "", "end date to query a chart from; e.g. 2022-08-31")
	viper.BindPFlags(runCmd.PersistentFlags())

}
