package functions

import (
	"os/exec"
	"fmt"
	"strconv"
)

//Requires xbacklight (Package: xorg-xbacklight)
func GetBrightness() string {
	cmd := 	exec.Command("xbacklight")
	brightness, err := cmd.Output()
	if err != nil {
		return ""
	}

	brightnessString := string(brightness[:len(brightness)-1])

	brightnessFloat, err := strconv.ParseFloat(brightnessString, 10)
	if err != nil {
		return ""
	}

	var formattedString string

	if brightnessFloat < 50.0 {
		formattedString = fmt.Sprintf("🔅 %0.00f", brightnessFloat)
	} else {
		formattedString = fmt.Sprintf("🔆 %0.00f", brightnessFloat)
	}
	

	return formattedString
}
