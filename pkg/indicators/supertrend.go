package indicators

import (
	"goblin-trader/pkg/common"
	"math"

	"github.com/sdcoffey/techan"
	log "github.com/sirupsen/logrus"
)

func SuperTrend(series *techan.TimeSeries, lookback int, multiplier float64) (ub, lb []float64, superTrend []bool) {
	// closePrices is a list[]float64 of closes for all datapoints in series
	closePrices := close(series)

	//TODO: check if length of series is longer than lookback
	// if series.Len() == 0 || series.Len() <= lookback {
	// 	log.Warn("The series must have at least one element and its length must be longer than the lookback period.")
	// 	return
	// }

	// get Average True Range
	a := techan.NewAverageTrueRangeIndicator(series, lookback)
	atr := common.Dump(a)
	hla := highLowAvg(series)
	ub, lb = calculateBands(hla, atr, multiplier)
	log.Debugf("ub: %v || lb: %v ", ub, lb)

	var prev int
	length := len(closePrices)
	superTrend = make([]bool, length)
	// Initialize the first SuperTrend value
	superTrend[0] = true
	// calculate supertrend values
	for i := 1; i < length; i++ {
		if closePrices[i] > ub[prev] {
			superTrend[i] = true
			lb[i] = math.NaN()
		} else if closePrices[i] < lb[prev] {
			superTrend[i] = false
			ub[i] = math.NaN()
		} else {
			superTrend[i] = superTrend[prev]

			if superTrend[i] {
				lb[i] = math.Max(lb[i], lb[prev])
				ub[i] = math.NaN()
			} else {
				ub[i] = math.Min(ub[i], ub[prev])
				lb[i] = math.NaN()
			}
		}
	}

	log.Debugf("final supertrend values: %v", superTrend)
	log.Debugf("final upperband values: %v", ub)
	log.Debugf("final lowerband values: %v", lb)

	return
}

func highLowAvg(series *techan.TimeSeries) []float64 {
	var hla []float64
	for _, s := range series.Candles {
		hla = append(hla, (s.MaxPrice.Add(s.MinPrice).Float() / 2))
	}
	return hla
}

func calculateBands(hla, atr []float64, multiplier float64) (ub, lb []float64) {
	ub = make([]float64, len(hla))
	lb = make([]float64, len(hla))

	for i := range hla {
		ub[i] = hla[i] + (multiplier * atr[i])
		lb[i] = hla[i] - (multiplier * atr[i])
	}

	return ub, lb
}

func close(series *techan.TimeSeries) []float64 {
	var c []float64
	for _, s := range series.Candles {
		c = append(c, s.ClosePrice.Float())
	}
	return c
}
