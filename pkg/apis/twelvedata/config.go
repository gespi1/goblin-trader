package twelvedata

type Config struct {
	Asset      string
	BaseURL    string `default:"https://api.twelvedata.com"`
	DateFormat string
	Interval   string
	Exchange   string
	EndDate    string
	StartDate  string
	Token      string
	Uri        Uri
}

type Uri struct {
	TimeSeries string `default:"/time_series?apikey={{ .Token }}&symbol={{ .Asset }}&interval={{ .Interval }}&exchange={{ .Exchange }}&start_date={{ .StartDate }}&end_date={{ .EndDate }}"`
}
