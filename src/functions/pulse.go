package functions

import (
	"github.com/vially/volumectl/pulseaudio"
	"fmt"
)


// GetPulseVolume returns the value of the volume level from pulseaudio
func GetPulseVolume(mutedFormat string, volumeFormat string) string {

	client := pulseaudio.New()
	volume := client.Volume
	muted := client.Muted

	if muted {
		return mutedFormat
	}

	volumeOutput := ""

	for i := 0; i < len(volumeFormat); i++ {
		char := string(volumeFormat[i])

		if char == "@" {
			i++
			nextChar := string(volumeFormat[i])

			if nextChar == "v" {
				volumeOutput += fmt.Sprintf("%v", volume)
			}	
		} else {
			volumeOutput += char
		}
	}
	
	return volumeOutput

}