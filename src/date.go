package main

import (
	"strings"
	"time"
)

//GetDate provides the current date in the format that the user specifies
func GetDate(config *configInterface) string {
	fullTime := time.Now()
	format := config.DateFormat

	if format == "" {
		format = "Jan 02 2006"
		return fullTime.Format(format)
	}

	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "mmm", "Jan")
	format = strings.ReplaceAll(format, "ddd", "Mon")

	formattedDate := fullTime.Format(format)

	return formattedDate
}
