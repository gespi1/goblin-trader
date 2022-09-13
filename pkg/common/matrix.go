package common

import (
	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/gonum/mat"
)

type Matrix struct {
	dataframe.DataFrame
}

func (m Matrix) At(i, j int) float64 {
	return m.Elem(i, j).Float()
}

func (m Matrix) T() mat.Matrix {
	return mat.Transpose{m}
}
