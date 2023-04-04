package common

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func CreateDataFrame(data ...interface{}) dataframe.DataFrame {
	if len(data)%2 != 0 {
		panic("createDataFrame expects an even number of arguments")
	}

	columns := make([]series.Series, len(data)/2)

	for i := 0; i < len(data); i += 2 {
		label, ok := data[i].(string)
		if !ok {
			panic("createDataFrame expects a string label for each column")
		}

		switch values := data[i+1].(type) {
		case []string:
			columns[i/2] = series.New(values, series.String, label)
		case []float64:
			columns[i/2] = series.New(values, series.Float, label)
		case []int:
			columns[i/2] = series.New(values, series.Int, label)
		case []bool:
			columns[i/2] = series.New(values, series.Bool, label)
		default:
			panic("createDataFrame: unsupported data type")
		}
	}

	return dataframe.New(columns...)
}

func WriteDFToFile(df dataframe.DataFrame, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := make([]string, df.Ncol())
	for j, columnName := range df.Names() {
		header[j] = fmt.Sprintf("%-15s", columnName)
	}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	// Write the data rows
	for i := 0; i < df.Nrow(); i++ {
		row := make([]string, df.Ncol())
		for j := 0; j < df.Ncol(); j++ {
			elem := df.Col(df.Names()[j]).Elem(i)
			if elem.Type() == series.Float {
				row[j] = fmt.Sprintf("%-15.2f", elem.Float())
			} else {
				row[j] = fmt.Sprintf("%-15s", elem.String())
			}
		}
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
