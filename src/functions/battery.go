package functions

import (
	"os/exec"
	"strings"
	"text/scanner"
	// "strconv"
)

//GetBatteryPercentage requires the acpi package https://sourceforge.net/projects/acpiclient/
func GetBatteryPercentage(batteryFormat string) string {

	var percentageString string

	cmd := exec.Command("acpi")

	output, err := cmd.Output()

	if err != nil {
		return ""
	}


	var s scanner.Scanner
	s.Init(strings.NewReader(string(output)))

	for i := 0; i < 6; i++ {
		s.Scan()
		if i == 5  {
			percentageString = s.TokenText()
		}
	}

	// percentage, err := strconv.ParseInt(percentageString, 10, 32)
	// if err != nil {
	// 	return ""
	// }


	formattedOutput := ""
	
	for i := 0; i < len(batteryFormat); i++ {
		char := string(batteryFormat[i])

		if char == "@" {
			nextChar := string(batteryFormat[i + 1])
			i++
			if nextChar == "b" {
				formattedOutput += percentageString
			}
		} else if char == "%"{
			formattedOutput += "%"
		} else {
			formattedOutput += string(char)
		}
	}

	// return fmt.Sprintf("⚡ %v%%", percentageString)
	return formattedOutput
}


