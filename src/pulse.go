package main

import (
	"fmt"
	"strings"

	"github.com/vially/volumectl/pulseaudio"
)

// GetPulseVolume returns the value of the volume level from pulseaudio
func GetPulseVolume(config *configInterface) string {

	mutedFormat := config.PulseMutedFormat
	volumeFormat := config.PulseVolumeFormat

	client := pulseaudio.New()
	volume := client.Volume
	muted := client.Muted

	if muted {
		return mutedFormat
	}

	// volumeOutput := ""

	// for i := 0; i < len(volumeFormat); i++ {
	// 	char := string(volumeFormat[i])

	// 	if char == "@" {
	// 		i++
	// 		nextChar := string(volumeFormat[i])

	// 		if nextChar == "v" {
	// 			volumeOutput += fmt.Sprintf("%v", volume)
	// 		}
	// 	} else {
	// 		volumeOutput += char
	// 	}
	// }

	return strings.ReplaceAll(volumeFormat, "@v", fmt.Sprintf("%v", volume))
}
