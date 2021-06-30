package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

//Requires xbacklight (Package: xorg-xbacklight)
func GetBrightness(brightnessChan chan string, config *configInterface) {
	for {
		cmd := exec.Command("xbacklight")
		brightness, err := cmd.Output()
		if err != nil {
			brightnessChan <- ""
			time.Sleep(time.Second)
			continue

		}

		brightnessString := string(brightness[:len(brightness)-1])

		brightnessFloat, err := strconv.ParseFloat(brightnessString, 10)
		if err != nil {
			brightnessChan <- ""
			time.Sleep(time.Second)
			continue
		}

		var formattedString string

		if brightnessFloat < 50.0 {
			formattedString = fmt.Sprintf("🔅 %v", int64(brightnessFloat))
		} else {
			formattedString = fmt.Sprintf("🔆 %v", int64(brightnessFloat))
		}

		brightnessChan <- formattedString
		time.Sleep(time.Second)
	}
}
