package main

import (
	"fmt"
	"strings"
	"time"

	// "strconv"
	"github.com/distatus/battery"
)

//GetBatteryPercentage requires the acpi package https://sourceforge.net/projects/acpiclient/
func GetBatteryPercentage(batteryChan chan string, config *configInterface) {
	for {
		battery, err := battery.Get(0)
		if err != nil {
			fmt.Println("Could not get battery info!")
			time.Sleep(time.Second)
			continue
		}

		percentageString := fmt.Sprintf("%.f", 100*battery.Current/battery.Full)

		formattedOutput := config.BatteryFormat

		formattedOutput = strings.ReplaceAll(formattedOutput, "@b", percentageString)

		batteryChan <- formattedOutput
		time.Sleep(10 * time.Second)
		continue
	}
}
