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
	"strings"
	"time"

	"goblin-trader/pkg/common"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	log "github.com/sirupsen/logrus"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
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

type XY struct{ X, Y float64 }

func (c *Config) TimeSeries() (*techan.TimeSeries, error) {
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
		period := techan.NewTimePeriod(time.Unix(start, 0), time.Hour*168)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(fmt.Sprintf("%v", v["open"]))
		candle.ClosePrice = big.NewFromString(fmt.Sprintf("%v", v["close"]))
		candle.MaxPrice = big.NewFromString(fmt.Sprintf("%v", v["high"]))
		candle.MinPrice = big.NewFromString(fmt.Sprintf("%v", v["low"]))

		if !series.AddCandle(candle) {
			log.Errorf("wasn't able to append candle %v", candle)
		}
	}

	return series, nil
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
