package common

import "github.com/sdcoffey/techan"

func GetPricesSlice(series *techan.TimeSeries) (prices []float64) {
	for _, s := range series.Candles {
		prices = append(prices, s.ClosePrice.Float())
	}
	return prices
}

func GetDatesInUnixSlice(series *techan.TimeSeries) (dates []float64) {
	for _, s := range series.Candles {
		dates = append(dates, float64(s.Period.End.Unix()))
	}
	return dates
}

func GetDatesHumanReadableSlice(series *techan.TimeSeries) (dates []string) {
	for _, s := range series.Candles {
		dates = append(dates, s.Period.End.String())
	}
	return dates
}
