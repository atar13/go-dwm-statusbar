package main 

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
	

	formattedOutput := batteryFormat
	
	formattedOutput = strings.ReplaceAll(formattedOutput, "@b", percentageString)

	return formattedOutput
}


