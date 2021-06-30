package main

import (
	"fmt"
	"mrogalski.eu/go/xbacklight"
	"time"
)

func GetBrightness(brightnessChan chan string, config *configInterface) {
	for {
		backlighter, err := xbacklight.NewBacklighterPrimaryScreen()
		if err != nil {
			fmt.Println("Error with connecting to X display:", err)
			brightnessChan <- ""
			time.Sleep(time.Second)
			continue

		}
		brightnessFloat, err := backlighter.Get()
		if err != nil {
			fmt.Println("Couldn't query backlight:", err)
			brightnessChan <- ""
			time.Sleep(time.Second)
			continue
		}
		brightnessFloat *= 100

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
