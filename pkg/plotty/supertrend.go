package plotty

import (
	"goblin-trader/pkg/common"
	"image/color"
	"math"

	"github.com/sdcoffey/techan"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type IndexedXYs []IndexedXY

func (ixys IndexedXYs) Len() int {
	return len(ixys)
}

func (ixys IndexedXYs) XY(i int) (float64, float64) {
	return ixys[i].X, ixys[i].Y
}

type IndexedXY struct {
	plotter.XY
	OriginalIndex int
}

// ub = upperband, lb = lowerband
func SuperTrend(series *techan.TimeSeries, ub, lb []float64, superTrend []bool) {
	// plotting
	var pricesY []float64
	var datesX []float64
	var datesHumanReadable []string

	pricesY = common.GetPricesSlice(series)
	datesX = common.GetDatesInUnixSlice(series)
	datesHumanReadable = common.GetDatesHumanReadableSlice(series)

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
		upperLine.Color = color.RGBA{R: 255, A: 255} // Red color
		p.Add(upperLine)
	}

	for _, segment := range lowerBandSegments {
		lowerLine, err := plotter.NewLine(common.XYPoints(segment))
		if err != nil {
			log.Fatalf("Failed to create lowerBand line: %v", err)
		}
		lowerLine.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green color
		p.Add(lowerLine)
	}

	line, err := plotter.NewLine(pricePlot)
	if err != nil {
		log.Fatalf("Failed to create new line: %v", err)
	}

	// Set line properties
	line.Color = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	line.Width = vg.Points(1) // line thickness

	// Add the price line to the plot
	p.Add(line)

	combinedSuperTrend := make(IndexedXYs, 0, len(superTrend))
	for i := range superTrend {
		if shouldPlotPoint(i, superTrend) {
			combinedSuperTrend = append(combinedSuperTrend, IndexedXY{XY: plotter.XY{X: float64(datesX[i]), Y: pricesY[i]}, OriginalIndex: i})
		}
	}

	scatterCombined, err := plotter.NewScatter(combinedSuperTrend)
	if err != nil {
		log.Fatalf("Failed to create SuperTrend scatter plot for combined values: %v", err)
	}
	scatterCombined.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		gs := draw.GlyphStyle{Radius: vg.Points(4), Shape: draw.CircleGlyph{}}
		originalIndex := combinedSuperTrend[i].OriginalIndex
		if superTrend[originalIndex] {
			gs.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green color
		} else {
			gs.Color = color.RGBA{R: 255, A: 255} // Red color
		}
		return gs
	}
	p.Add(scatterCombined)

	if err := p.Save(15*vg.Inch, 10*vg.Inch, "supertrend.png"); err != nil {
		log.Fatalf("Failed to save plot: %v", err)
	}

	df := common.CreateDataFrame("Date", datesHumanReadable, "Price", pricesY, "SuperTrend", superTrend, "LowerBand", lb, "UpperBand", ub)
	common.WriteDFToFile(df, "./superTrendDF.csv")
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

func shouldPlotPoint(index int, superTrend []bool) bool {
	if index == 0 {
		return true
	}
	return superTrend[index] != superTrend[index-1]
}
