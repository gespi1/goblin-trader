package common

import (
	"time"

	"gonum.org/v1/plot"
)

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
