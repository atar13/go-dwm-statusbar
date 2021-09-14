package main

import (
	"fmt"
	"os"
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
		switch statusString {
		case "Discharging\n":
			statusOutput = config.DischargingIndicator
		case "Charging\n":
			statusOutput = config.ChargingIndicator
		default:

		}

		formattedOutput := config.BatteryFormat

		formattedOutput = strings.ReplaceAll(formattedOutput, "@b", percentageString)
		formattedOutput = statusOutput + formattedOutput

		return formattedOutput
	}
	return ""
}
