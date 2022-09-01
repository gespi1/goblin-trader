package common

import (
	"time"
)

const YYYYMMDD = "2006-01-02"

// returns currentDate and Date from a Year ago
func getYearDates() (string, string) {

	currentDate := time.Now()
	today := currentDate.Format(YYYYMMDD)

	minusYearDate := currentDate.AddDate(-1, 0, 0)
	aYearAgo := minusYearDate.Format(YYYYMMDD)

	return aYearAgo, today
}

// determines if dates are empty and if they are return YTD dates
func DetermineDate(startDate string, endDate string) (string, string) {
	if startDate == "" && endDate == "" {
		return getYearDates()
	} else if startDate != "" && endDate == "" {
		_, e := getYearDates()
		return startDate, e
	} else {
		return startDate, endDate
	}
}
