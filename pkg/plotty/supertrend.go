package plotty

import (
	"image/color"
	"math"
	"time"

	"github.com/sdcoffey/techan"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// ub = upperband, lb = lowerband
func SuperTrend(series *techan.TimeSeries, ub, lb []float64, superTrend []bool) {
	// plotting
	var pricesY []float64
	var datesX []float64

	for _, s := range series.Candles {
		pricesY = append(pricesY, s.ClosePrice.Float())
		datesX = append(datesX, float64(s.Period.End.Unix()))
	}

	p := plot.New()
	// create labeles; title, x axis label, y axis label
	p.Title.Text = "Prices, SuperTrend Bands and SuperTrend"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"
	p.Add(plotter.NewGrid())

	// Set the custom time ticker for the x-axis
	p.X.Tick.Marker = TimeTicker{NLabels: 12}
	p.Y.Tick.Marker = customTickMarker{NumTicks: 20}

	pricePlot := make(plotter.XYs, len(pricesY))
	superTrendPlot := make(plotter.XYs, 0, len(superTrend))

	upperBandSegments := createLineSegments(datesX, ub)
	lowerBandSegments := createLineSegments(datesX, lb)

	for i := range pricesY {
		pricePlot[i].X = float64(datesX[i])
		pricePlot[i].Y = pricesY[i]

		if superTrend[i] {
			superTrendPlot = append(superTrendPlot, plotter.XY{X: float64(datesX[i]), Y: pricesY[i]})
		}
	}

	for _, segment := range upperBandSegments {
		upperLine, err := plotter.NewLine(XYPoints(segment))
		if err != nil {
			log.Fatalf("Failed to create upperBand line: %v", err)
		}
		upperLine.Color = color.RGBA{R: 255, A: 255}
		p.Add(upperLine)
	}

	for _, segment := range lowerBandSegments {
		lowerLine, err := plotter.NewLine(XYPoints(segment))
		if err != nil {
			log.Fatalf("Failed to create lowerBand line: %v", err)
		}
		lowerLine.Color = color.RGBA{B: 255, A: 255}
		p.Add(lowerLine)
	}

	// fmt.Println("1")
	// fmt.Println(pricePlot)
	// fmt.Println("2")
	// fmt.Println(upperBandPlot)
	// fmt.Println("3")
	// fmt.Println(lowerBandPlot)

	err := plotutil.AddLines(p,
		"Price", pricePlot,
	)
	if err != nil {
		log.Fatalf("Failed to add lines to plot: %v", err)
	}

	scatter, err := plotter.NewScatter(superTrendPlot)
	if err != nil {
		log.Fatalf("Failed to create SuperTrend scatter plot: %v", err)
	}
	scatter.GlyphStyle.Radius = vg.Points(2)
	scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green color
	p.Add(scatter)

	if err := p.Save(10*vg.Inch, 5*vg.Inch, "supertrend.png"); err != nil {
		log.Fatalf("Failed to save plot: %v", err)
	}
}

type TimeTicker struct {
	// How many labels to show on the x-axis
	NLabels int
}

func (t TimeTicker) Ticks(min, max float64) []plot.Tick {
	delta := (max - min) / float64(t.NLabels)

	ticks := make([]plot.Tick, t.NLabels)
	for i := 0; i < t.NLabels; i++ {
		value := min + delta*float64(i)
		ticks[i] = plot.Tick{
			Value: value,
			Label: formatEpochTime(int64(value)),
		}
	}
	return ticks
}

func formatEpochTime(epoch int64) string {
	t := time.Unix(epoch, 0)
	return t.Format("2006-01-02")
}

func createLineSegments(epochTimes, data []float64) [][]plotter.XY {
	segments := make([][]plotter.XY, 0)
	currentSegment := make([]plotter.XY, 0)

	for i := range data {
		if !math.IsNaN(data[i]) {
			currentSegment = append(currentSegment, plotter.XY{X: float64(epochTimes[i]), Y: data[i]})
		} else {
			if len(currentSegment) > 0 {
				segments = append(segments, currentSegment)
				currentSegment = make([]plotter.XY, 0)
			}
		}
	}

	if len(currentSegment) > 0 {
		segments = append(segments, currentSegment)
	}

	return segments
}

type XYPoints []plotter.XY

func (p XYPoints) Len() int {
	return len(p)
}

func (p XYPoints) XY(i int) (x, y float64) {
	return p[i].X, p[i].Y
}

type customTickMarker struct {
	plot.DefaultTicks
	NumTicks int
}

func (t customTickMarker) Ticks(min, max float64) []plot.Tick {
	ticks := t.DefaultTicks.Ticks(min, max)
	if t.NumTicks <= 0 {
		return ticks
	}

	newTicks := make([]plot.Tick, 0, t.NumTicks)
	step := len(ticks) / t.NumTicks
	if step == 0 {
		step = 1
	}

	for i := 0; i < len(ticks); i += step {
		newTicks = append(newTicks, ticks[i])
	}
	return newTicks
}
