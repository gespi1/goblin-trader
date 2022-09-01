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

	"github.com/go-gota/gota/dataframe"
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
	Meta   map[string]string  `json:"meta"`
	Values []TimeSeriesValues `json:"values"`
	Status string             `json:"status"`
}

type TimeSeriesValues struct {
	Datetime string  `json:"datetime"`
	Open     string  `json:"open"`
	High     string  `json:"high"`
	Low      string  `json:"low"`
	Close    string  `json:"close"`
	Unixtime float64 `json:"unixtime"`
}

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

	err = json.Unmarshal(body, &ts)

	if err != nil {
		log.Errorf("wasnt able to unmarshal response; %v", err)
	}

	c.dateTimeToUnix(ts.Values)
	df := dataframe.LoadStructs(ts.Values)
	// fmt.Println(df.Select([]string{"Datetime", "Close"}))
	// fmt.Println(df.Col("Close"))
	fmt.Println(df.Select([]string{"Unixtime", "Close"}))

	XYs := makeXYs(df.Select([]string{"Unixtime", "Close"}).Records())

	// plotting
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

func (c *Config) dateTimeToUnix(datetimes []TimeSeriesValues) {
	for i, d := range datetimes {

		tm, err := time.Parse(c.DateFormat, d.Datetime)
		if err != nil {
			log.Errorf("wasn't able to parse time %v: %v", d.Datetime, err)
		}

		datetimes[i].Unixtime = float64(tm.Unix())
	}
}

func determineTimeFormat(interval string) string {
	match, err := regexp.MatchString(`^(\d{1}|\d{2})h$`, interval)
	if err != nil {
		fmt.Errorf("matching regex %v: %v", interval, err)
	}

	if match || strings.Contains(interval, "min") {
		fmt.Println("h or min")
		return "2006-01-02 15:04:05"
	} else {
		fmt.Println("NOT h or min")
		return "2006-01-02"
	}

}

func makeXYs(data [][]string) plotter.XYs {
	var XYs plotter.XYs

	// remove the titles from [][]string created by the dataframe
	data = append(data[:0], data[0+1:]...)

	for _, item := range data {
		x64, err := strconv.ParseFloat(item[0], 64)
		if err != nil {
			log.Errorf("wasn't able to parse x64 %v: %v", item[0], err)
		}
		y64, err := strconv.ParseFloat(item[1], 64)
		if err != nil {
			log.Errorf("wasn't able to parse y64 %v: %v", item[1], err)
		}

		XYs = append(XYs, plotter.XY{X: x64, Y: y64})
	}
	return XYs
}
