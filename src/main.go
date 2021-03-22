package main

import (
	"fmt"

	// "io"
	"io/ioutil"
	"os"
	"os/exec"

	"strconv"
	"time"

	// F "functions"
	"gopkg.in/yaml.v2"
)


type configInterface struct{
	Modules			[]string	`yaml:"Modules"`
	ModuleSeperator	string		`yaml:"ModuleSeperator"`
	RefreshConfig	bool 		`yaml:"RefreshConfig"`
	RefreshConfigRate string	`yaml:"RefreshConfigRate"`
	TimeFormat 		string 		`yaml:"TimeFormat"`
	TwentyFourHour 	bool 		`yaml:"TwentyFourHour"`
	DateFormat		string		`yaml:"DateFormat"`
	PlayingFormat	string		`yaml:"PlayingFormat"`
	PausedFormat	string		`yaml:"PausedFormat"`
	MprisMaxLength 	string		`yaml:"MprisMaxLength"`
	ScrollMpris		bool		`yaml:"ScrollMpris"`
	CPUTempUnits 	string 		`yaml:"CPUTempUnits"`
	BatteryFormat 	string 		`yaml:"BatteryFormat"`
	RAMDisplay 		string		`yaml:"RAMDisplay"`
	RAMRawUnit		string		`yaml:"RAMRawUnit"`
	RAMRawFormat	string 		`yaml:"RAMRawFormat"`
	PulseMutedFormat string 	`yaml:"PulseMutedFormat"`
	PulseVolumeFormat string 	`yaml:"PulseVolumeFormat"`
}



func main()  {

    var config configInterface
	parsedConfig := config.retrieveConfig()

	fmt.Println("Config file successfully loaded")
	fmt.Println(parsedConfig)

	loopCounter := 0

	for {
		loopCounter++
		modules := parsedConfig.Modules

		configRefreshRate, err := strconv.Atoi(parsedConfig.RefreshConfigRate)
		if err != nil {
			configRefreshRate = 100
		}


		output := ""
		for idx, module := range modules {
			moduleData := ""
			switch module {
				case "time":
					moduleData += GetTime(parsedConfig.TimeFormat, parsedConfig.TwentyFourHour)
				case "date":
					moduleData += GetDate(parsedConfig.DateFormat)
				case "mpris":
					moduleData += GetMpris(parsedConfig.PlayingFormat, parsedConfig.PausedFormat, parsedConfig.MprisMaxLength, parsedConfig.ScrollMpris)
				case "cpu":
					moduleData += GetCPUTemp(parsedConfig.CPUTempUnits)
					// moduleData += F.GetCPUUsage()
				case "battery":
					moduleData += GetBatteryPercentage(parsedConfig.BatteryFormat)
				case "ram":
					if parsedConfig.RAMDisplay == "Percentage" {
						moduleData += GetRAMUsage("format placeholder", false)
					} else if parsedConfig.RAMDisplay == "Raw" {
						moduleData += GetRAMData(parsedConfig.RAMRawFormat, parsedConfig.RAMRawUnit)
					}
				case "brightness":
					moduleData += GetBrightness()
				case "pulse":
					moduleData += GetPulseVolume(parsedConfig.PulseMutedFormat, parsedConfig.PulseVolumeFormat)
			}
			if moduleData == "" {
				continue
			}

			if idx != len(modules) - 1 {
				moduleData += parsedConfig.ModuleSeperator
			}

			output += moduleData
		}


		xsetroot := "xsetroot"
		arg1 := "-name"

		cmd := exec.Command(xsetroot, arg1, output)

		_, err = cmd.Output()

		if err != nil {
			fmt.Println(err)
			return 
		}
		if parsedConfig.RefreshConfig && loopCounter % configRefreshRate == 0 {
			parsedConfigChan := make(chan configInterface)	
			go refreshConfig(parsedConfigChan, parsedConfig) 
			parsedConfig = <- parsedConfigChan
			fmt.Println(parsedConfig)
		}
		time.Sleep(100 * time.Millisecond)
	}

}



func (config *configInterface) retrieveConfig() configInterface {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/go-dwm-statusbar/config.yaml", os.Getenv("HOME")))

	fmt.Println("Updating config")
	if err != nil {
		fmt.Println("Error with reading config file")
    }

	yaml.Unmarshal(data, config)

	data = nil
	
	return *config
}

func refreshConfig(config chan configInterface, configStruct configInterface) {
		config <- configStruct.retrieveConfig()
}
