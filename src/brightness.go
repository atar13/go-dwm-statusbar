package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetBrightness(config *configInterface) string {
	rootPath := "/sys/class/backlight"
	backlightSources, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, backlightSource := range backlightSources {

		brightnessFile, err := os.ReadFile(strings.Join([]string{rootPath, backlightSource.Name(), "brightness"}, "/"))
		if err != nil {
			fmt.Println("Error reading brightness")
			return ""
		}
		brightnessString := string(brightnessFile)
		brightnessString = strings.Replace(brightnessString, "\n", "", -1)
		brightnessInt, err := strconv.Atoi(brightnessString)
		if err != nil {
			fmt.Println("Error parsing brightness")
			return ""
		}

		maxBrightnessFile, err := os.ReadFile(strings.Join([]string{rootPath, backlightSource.Name(), "max_brightness"}, "/"))
		if err != nil {
			fmt.Println("Error reading brightness")
			return ""
		}
		maxBrightnessString := string(maxBrightnessFile)
		maxBrightnessString = strings.Replace(maxBrightnessString, "\n", "", -1)
		maxBrightnessInt, err := strconv.Atoi(maxBrightnessString)
		if err != nil {
			fmt.Println("Error parsing brightness")
			return ""
		}

		brightnessFloat := (float32(brightnessInt) / float32(maxBrightnessInt)) * 100
		fmt.Println(brightnessFloat)

		var formattedString string

		if brightnessFloat < 50.0 {
			formattedString = fmt.Sprintf("🔅 %d", int64(brightnessFloat))
		} else {
			formattedString = fmt.Sprintf("🔆 %d", int64(brightnessFloat))
		}

		return formattedString
	}
	return ""
}
