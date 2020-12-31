package functions

import (
	"github.com/vially/volumectl/pulseaudio"
	"fmt"
)

func GetPulseVolume() string {

	client := pulseaudio.New()
	volume := client.Volume
	return fmt.Sprint(volume)
}