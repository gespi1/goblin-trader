package strategies

func SupertrendStrategy(prices []float64, supertrend []bool) []Trade {
	trades := []Trade{}

	// Loop through the price data and supertrend signals
	for i := 1; i < len(prices); i++ {
		prevSignal, currSignal := supertrend[i-1], supertrend[i]

		// Check for buy signal: previous supertrend was false (bearish), and current supertrend is true (bullish)
		if !prevSignal && currSignal {
			trades = append(trades, Trade{
				EntryPrice: prices[i],
				IsLong:     true,
			})
		}

		// Check for sell signal: previous supertrend was true (bullish), and current supertrend is false (bearish)
		if prevSignal && !currSignal {
			trades = append(trades, Trade{
				EntryPrice: prices[i],
				IsLong:     false,
			})
		}
	}

	return trades
}

func (r *TradingRules) SuperTrendExecuteTrades(prices []float64, trades []Trade) ([]Trade, []float64, float64) {
	executedTrades := []Trade{}
	tradeProfits := []float64{}

	positionOpen := false
	var entryPrice float64
	var currentPosition Trade

	// Set starting capital
	capital := r.StartingCapital

	for i := 0; i < len(prices); i++ {
		currentPrice := prices[i]

		for _, trade := range trades {
			if trade.EntryPrice == currentPrice {
				if positionOpen {
					// Close the current position and open a new one in the opposite direction
					profit := calculateProfit(entryPrice, currentPrice, currentPosition.IsLong)
					tradeProfits = append(tradeProfits, profit)

					// Update capital
					capital += capital * profit
				}

				positionOpen = true
				entryPrice = currentPrice
				currentPosition = trade
				break
			}
		}

		if positionOpen {
			// Calculate stop-loss and take-profit levels
			stopLoss := entryPrice * (1.0 - r.StopLossPct)
			takeProfit := entryPrice * (1.0 + r.TakeProfitPct)

			var isLong bool
			for _, trade := range trades {
				if trade.EntryPrice == entryPrice {
					isLong = trade.IsLong
					break
				}
			}

			if isLong {
				if currentPrice <= stopLoss || currentPrice >= takeProfit {
					// Close the position
					profit := calculateProfit(entryPrice, currentPrice, isLong)
					tradeProfits = append(tradeProfits, profit)

					// Update capital
					capital += capital * profit
					positionOpen = false
				}
			} else {
				if currentPrice >= stopLoss || currentPrice <= takeProfit {
					// Close the position
					profit := calculateProfit(entryPrice, currentPrice, isLong)
					tradeProfits = append(tradeProfits, profit)

					// Update capital
					capital += capital * profit
					positionOpen = false
				}
			}
		}
	}

	return executedTrades, tradeProfits, capital
}

func calculateProfit(entryPrice, exitPrice float64, isLong bool) float64 {
	if isLong {
		return (exitPrice - entryPrice) / entryPrice
	} else {
		return (entryPrice - exitPrice) / entryPrice
	}
}
