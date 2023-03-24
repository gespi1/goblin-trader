package common

import "gonum.org/v1/plot/plotter"

type XYPoints []plotter.XY

func (p XYPoints) Len() int {
	return len(p)
}

func (p XYPoints) XY(i int) (x, y float64) {
	return p[i].X, p[i].Y
}
