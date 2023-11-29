## Goblin Trader: The GoLang Trading Bot
Opensource Trading Bot based on Go. Python would probably be better suited for this but eh. 

The goal of this project is to have a bot that trades for me during the work day in a reliable and stress free way. Stress free you may ask, well yes. With the combination of the SuperTrend and using longer time windows of an Asset we are able to get reliable Buy and Sell signals. 

This project also supports short term traders, so no worries for the degens out there. 



So far this proven to work with Cryptocurrency key pairs e.g. `BTC/USD`, `ETH/USD`. Check out other supported [key pairs](https://api.twelvedata.com/cryptocurrencies).


### Roadmap
 - Gather Datapoints from TwelveData API ✅
 - Transform TwelveData datapoints into DataFrame ✅
 - Code SuperTrend Formula to use against DataFrame ✅
 - Display SuperTrend and Asset price on a graph using Plotty and render it as a PNG ✅
 - Represent a SuperTrend against an Asset via CSV ✅
 - Create Backtesting
 - Create Trading Strategy
 - Add support for exhanges (Bianance, Coinbase and others)
 - Have the bot trade for me!

 ### Usage so far
 ```
 export TWELVE_TOKEN=<KEY>
 go run main.go chart -a BTC/USD  -t 1day  -s 2023-01-01
 ```

 will result in a a chart like this:
![BTC_USD since Jan 2023](images/btc_usd-1day-Jan_2023.png)``