package main 

import (
	"time"
	"strings"
	)
// import "fmt"

const defaultFormat = "03:04:05"
const default24Format = "15:04:05"

//GetTime retrieves the current time and formats it according to the user's configuration
func GetTime(format string, twentyFourHour bool) string{
	fullTime := time.Now()

	if format == "" {
		if twentyFourHour {
			return fullTime.Format(default24Format)
		}
			return fullTime.Format(defaultFormat)
	}

	format = strings.ReplaceAll(format, "MM", "04")
	format = strings.ReplaceAll(format, "SS", "05")

	var formattedTime string
	if twentyFourHour {
		format = strings.ReplaceAll(format, "HH", "15")
		formattedTime = fullTime.Format(format)
	} else {
		format = strings.ReplaceAll(format, "HH", "03")
		formattedTime = fullTime.Format(format)
	}
	return formattedTime
}


