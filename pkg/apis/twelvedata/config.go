package twelvedata

type Config struct {
	Asset     string `json:"asset"`
	BaseURL   string `json:"base_url" default:"https://api.twelvedata.com"`
	Interval  string `json:"interval"`
	Token     string `json:"token"`
	Exchange  string `json:"exchange"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Uri       Uri    `json:"uri"`
}

type Uri struct {
	TimeSeries string `json:"time_series" default:"/time_series?apikey={apikey}&symbol={symbol}&interval={interval}&exchange={exchange}&start_date={start_date}&end_date={end_date}"`
}
