package common

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func SetLogLevel(loglevel string) {
	l := strings.ToLower(loglevel)
	if l == "info" {
		log.SetLevel(log.InfoLevel)
	} else if l == "warn" {
		log.SetLevel(log.WarnLevel)
	} else if l == "error" {
		log.SetLevel(log.ErrorLevel)
	} else if l == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if l == "fatal" {
		log.SetLevel(log.FatalLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
		log.Warn("no log level matched setting to default, ERROR")
	}
}
