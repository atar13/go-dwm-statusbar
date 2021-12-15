package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetBatteryPercentage(config *configInterface) string {
	// looks through the first 10 battery devices for the first one with battery information
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/sys/class/power_supply/BAT%d/", i)
		capacityPath := path + "capacity"
		statusPath := path + "status"

		capacityFile, err := os.ReadFile(capacityPath)
		if err != nil {
			fmt.Printf("Error with reading battery #%d information\n", i)
			continue
		}
		statusFile, err := os.ReadFile(statusPath)
		if err != nil {
			fmt.Printf("Error with reading battery #%d information\n", i)
			continue
		}

		percentageString := string(capacityFile)
		percentageString = strings.Replace(percentageString, "\n", "", -1)

		statusString := string(statusFile)
		statusOutput := ""
		percentageInt, err := strconv.Atoi(percentageString)
		if err != nil {
			statusOutput = ""
		}
		switch statusString {
		// case "Discharging\n":
		// 	// statusOutput = config.DischargingIndicator
		// 	if percentageInt <= 10 {
		// 		statusOutput = config.TenDischargingIndicator
		// 	} else if percentageInt <= 20 {
		// 		statusOutput = config.TwentyDischargingIndicator
		// 	} else if percentageInt <= 30 {
		// 		statusOutput = config.ThirtyDischargingIndicator
		// 	} else if percentageInt <= 40 {
		// 		statusOutput = config.FortyDischargingIndicator
		// 	} else if percentageInt <= 50 {
		// 		statusOutput = config.FiftyDischargingIndicator
		// 	} else if percentageInt <= 60 {
		// 		statusOutput = config.SixtyDischargingIndicator
		// 	} else if percentageInt <= 70 {
		// 		statusOutput = config.SeventyDischargingIndicator
		// 	} else if percentageInt <= 80 {
		// 		statusOutput = config.EightyDischargingIndicator
		// 	} else if percentageInt <= 90 {
		// 		statusOutput = config.NinetyDischargingIndicator
		// 	} else {
		// 		statusOutput = config.FullDischargingIndicator
		// 	}
		case "Charging\n":
			// statusOutput = config.ChargingIndicator
			if percentageInt <= 10 {
				statusOutput = config.TenChargingIndicator
			} else if percentageInt <= 20 {
				statusOutput = config.TwentyChargingIndicator
			} else if percentageInt <= 30 {
				statusOutput = config.ThirtyChargingIndicator
			} else if percentageInt <= 40 {
				statusOutput = config.FortyChargingIndicator
			} else if percentageInt <= 50 {
				statusOutput = config.FiftyChargingIndicator
			} else if percentageInt <= 60 {
				statusOutput = config.SixtyChargingIndicator
			} else if percentageInt <= 70 {
				statusOutput = config.SeventyChargingIndicator
			} else if percentageInt <= 80 {
				statusOutput = config.EightyChargingIndicator
			} else if percentageInt <= 90 {
				statusOutput = config.NinetyChargingIndicator
			} else {
				statusOutput = config.FullChargingIndicator
			}
		default:
			if percentageInt <= 10 {
				statusOutput = config.TenDischargingIndicator
			} else if percentageInt <= 20 {
				statusOutput = config.TwentyDischargingIndicator
			} else if percentageInt <= 30 {
				statusOutput = config.ThirtyDischargingIndicator
			} else if percentageInt <= 40 {
				statusOutput = config.FortyDischargingIndicator
			} else if percentageInt <= 50 {
				statusOutput = config.FiftyDischargingIndicator
			} else if percentageInt <= 60 {
				statusOutput = config.SixtyDischargingIndicator
			} else if percentageInt <= 70 {
				statusOutput = config.SeventyDischargingIndicator
			} else if percentageInt <= 80 {
				statusOutput = config.EightyDischargingIndicator
			} else if percentageInt <= 90 {
				statusOutput = config.NinetyDischargingIndicator
			} else {
				statusOutput = config.FullDischargingIndicator
			}

		}

		formattedOutput := config.BatteryFormat

		formattedOutput = strings.ReplaceAll(formattedOutput, "@b", percentageString)
		formattedOutput = statusOutput + formattedOutput

		return formattedOutput
	}
	return ""
}
