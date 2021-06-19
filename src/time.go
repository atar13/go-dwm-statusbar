package main 

import (
	"time"
	"strings"
	)
// import "fmt"

const defaultFormat = "03:04:05"
const default24Format = "15:04:05"

//GetTime retrieves the current time and formats it according to the user's configuration
func GetTime(timeChan chan string, config *configInterface) {
	for {
		fullTime := time.Now()
		format := config.TimeFormat
		twentyFourHour := config.TwentyFourHour

		if format == "" {
			if twentyFourHour {
				timeChan <- fullTime.Format(default24Format)
			} else {
				timeChan <- fullTime.Format(defaultFormat)
			}
			time.Sleep(time.Second)
			continue	
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
		// return formattedTime
		timeChan <- formattedTime
		time.Sleep(time.Second)
	}
}


