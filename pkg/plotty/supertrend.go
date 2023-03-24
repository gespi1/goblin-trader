package plotty

import (
	"goblin-trader/pkg/common"
	"image/color"
	"math"

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
	p.X.Tick.Marker = common.TimeTicker{NLabels: 12}

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
		upperLine, err := plotter.NewLine(common.XYPoints(segment))
		if err != nil {
			log.Fatalf("Failed to create upperBand line: %v", err)
		}
		upperLine.Color = color.RGBA{R: 255, A: 255}
		p.Add(upperLine)
	}

	for _, segment := range lowerBandSegments {
		lowerLine, err := plotter.NewLine(common.XYPoints(segment))
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
