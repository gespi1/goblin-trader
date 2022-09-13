package twelvedata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"goblin-trader/pkg/common"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	log "github.com/sirupsen/logrus"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func Init(v *viper.Viper) *Config {
	uri := Uri{}

	common.CheckForToken(v, "TWELVE_TOKEN")
	startDate, endDate := common.DetermineDate(v.GetString("start-date"), v.GetString("end-date"))

	tdConf := new(Config)
	tdConf.Asset = v.GetString("asset")
	tdConf.DateFormat = determineTimeFormat(v.GetString("interval"))
	tdConf.Interval = v.GetString("interval")
	tdConf.Exchange = v.GetString("exchange")
	tdConf.EndDate = endDate
	tdConf.StartDate = startDate
	tdConf.Token = v.GetString("TWELVE_TOKEN")
	tdConf.Uri = uri

	defaults.SetDefaults(tdConf)

	return tdConf
}

type TimeSeries struct {
	Meta   map[string]string        `json:"meta"`
	Values []map[string]interface{} `json:"values"`
	Status string                   `json:"status"`
}

// type TimeSeriesValues struct {
// 	Datetime string  `json:"datetime"`
// 	Open     string  `json:"open"`
// 	High     string  `json:"high"`
// 	Low      string  `json:"low"`
// 	Close    string  `json:"close"`
// 	Unixtime float64 `json:"unixtime"`
// 	Volume   float64 `json:"volume"`
// }

type XY struct{ X, Y float64 }

func (c *Config) TimeSeries() {
	var ts TimeSeries

	// construct URL
	url := c.BaseURL + c.Uri.TimeSeries
	log.Debugf("url %v", url)

	// interpolate URL
	var urlBytes bytes.Buffer
	t, _ := template.New("timeseries").Parse(url)
	t.Execute(&urlBytes, c)

	url = urlBytes.String()
	log.Debugf("url rendered %v", url)

	// setup http client
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("client could not create request: %s", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("client error making http request: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	log.Infof("client received response, status code: %v", res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("**BODY**: %v", string(body))

	err = json.Unmarshal(body, &ts)

	if err != nil {
		log.Errorf("wasnt able to unmarshal response; %v", err)
	}

	log.Debugf("Unmarshaled timeseries values: %v", ts.Values)

	series := techan.NewTimeSeries()

	for _, v := range ts.Values {
		start := c.dateTimeToUnix(fmt.Sprintf("%v", v["datetime"]))
		period := techan.NewTimePeriod(time.Unix(start, 0), time.Hour*24)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(fmt.Sprintf("%v", v["open"]))
		candle.ClosePrice = big.NewFromString(fmt.Sprintf("%v", v["close"]))
		candle.MaxPrice = big.NewFromString(fmt.Sprintf("%v", v["high"]))
		candle.MinPrice = big.NewFromString(fmt.Sprintf("%v", v["low"]))

		if !series.AddCandle(candle) {
			log.Errorf("wasn't able to append candle %v", candle)
		}
	}

	// plotting
	var closeY []string
	var dateX []string

	for _, s := range series.Candles {
		closeY = append(closeY, s.ClosePrice.String())
		dateX = append(dateX, strconv.FormatInt(s.Period.End.Unix(), 10))
	}

	XYs, err := makeXYs(dateX, closeY)
	if err != nil {
		log.Error(err)
	}

	f, err := os.Create("out.png")
	if err != nil {
		log.Errorf("could not create out.png: %v", err)
	}

	p := plot.New()
	line, err := plotter.NewLine(XYs)
	if err != nil {
		log.Error(err)
	}

	p.Add(line)

	wt, err := p.WriterTo(1024, 512, "png")
	if err != nil {
		log.Errorf("could not create writer: %v", err)
	}

	_, err = wt.WriteTo(f)
	if err != nil {
		log.Errorf("could not create out.png: %v", err)
	}

	if err := f.Close(); err != nil {
		log.Errorf("could not close out.png: %v", err)
	}

}

func (c *Config) dateTimeToUnix(datetime string) int64 {
	tm, err := time.Parse(c.DateFormat, datetime)
	if err != nil {
		log.Errorf("wasn't able to parse time %v: %v", datetime, err)
	}

	return tm.Unix()
}

func determineTimeFormat(interval string) string {
	match, err := regexp.MatchString(`^(\d{1}|\d{2})h$`, interval)
	if err != nil {
		log.Errorf("matching regex %v: %v", interval, err)
	}

	if match || strings.Contains(interval, "min") {
		log.Debug("Time format in 'h' or 'min': 2006-01-02 15:04:05")
		return "2006-01-02 15:04:05"
	} else {
		log.Debug("Time format NOT 'h' or 'min': 2006-01-02")
		return "2006-01-02"
	}

}

func makeXYs(x, y []string) (plotter.XYs, error) {
	var XYs plotter.XYs
	if len(x) == len(y) {
		log.Info("x and y match length. Creating X and Y points for graph")

		for i, _ := range x {
			fmt.Println(x[i])
			x64, err := strconv.ParseFloat(x[i], 64)
			if err != nil {
				log.Errorf("wasn't able to parse x64 %v: %v", x[0], err)
			}
			y64, err := strconv.ParseFloat(y[i], 64)
			if err != nil {
				log.Errorf("wasn't able to parse y64 %v: %v", y[1], err)
			}

			XYs = append(XYs, plotter.XY{X: x64, Y: y64})
		}
		return XYs, nil
	} else {
		return nil, fmt.Errorf("x:( %v )and y:( %v ) DON'T match length", x, y)
	}

}
