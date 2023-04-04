package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
		// if no startdate and enddate are defined get todays (date - 1year)
		return getYearDates()
	} else if startDate != "" && endDate == "" {
		// if startdate defined and not enddate get (today - startdate)
		_, today := getYearDates()
		return startDate, today
	} else {
		// if no ifs match just return startdate and enddate
		return startDate, endDate
	}
}

// converts datetime to unix
func DateTimeToUnix(dateformat, datetime string) int64 {
	tm, err := time.Parse(dateformat, datetime)
	if err != nil {
		log.Errorf("wasn't able to parse time %v: %v", datetime, err)
	}

	return tm.Unix()
}

// determines if we should use hours:minutes:seconds on our timestamp
func DetermineTimeFormat(interval string) string {
	match, err := regexp.MatchString(`^(\d{1}|\d{2})(h|min)$`, interval)
	if err != nil {
		log.Errorf("matching regex %v: %v", interval, err)
	}

	if match || strings.Contains(interval, "min") {
		// for hours and minutes we need the timestamp
		log.Debug("Time format in 'h' or 'min': 2006-01-02 15:04:05")
		return "2006-01-02 15:04:05"
	} else {
		// for non-hours and non-minutes we do not need the time stamp. Date is suffice
		log.Debug("Time format NOT 'h' or 'min': 2006-01-02")
		return "2006-01-02"
	}

}

// determines an interval (week,day,h,min) in minutes
func DetermineIntervalInMin(interval string) (int, error) {
	// Compile the regular expression.
	pattern, err := regexp.Compile(`^(\d{1,2})(day|min|h|month|week)$`)
	if err != nil {
		log.Errorf("Error: %v", err)
		return 0, err
	}

	// Find and capture the matching groups.
	matches := pattern.FindStringSubmatch(interval)

	if len(matches) > 0 {
		// The full match is always the first element, so the last group is at index len(matches)-1.
		lastGroup := matches[len(matches)-1]
		duration, _ := strconv.Atoi(matches[1])

		if lastGroup == "day" {
			return duration * 24 * 60, nil
		} else if lastGroup == "min" {
			return duration, nil
		} else if lastGroup == "h" {
			return duration * 60, nil
		} else if lastGroup == "month" {
			return duration * 30 * 24 * 60, nil
		} else if lastGroup == "week" {
			return duration * 7 * 24 * 60, nil
		} else {
			return 0, fmt.Errorf("unexpected match: %s", lastGroup)
		}

	} else {
		return 0, fmt.Errorf("the interval string does not match the regex pattern (day|min|h|month|week)")
	}
}
