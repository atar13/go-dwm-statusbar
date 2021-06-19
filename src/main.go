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
	MprisScrollSpeed string		`yaml:"MprisScrollSpeed"`
	CPUTempUnits 	string 		`yaml:"CPUTempUnits"`
	BatteryFormat 	string 		`yaml:"BatteryFormat"`
	RAMDisplay 		string		`yaml:"RAMDisplay"`
	RAMRawUnit		string		`yaml:"RAMRawUnit"`
	RAMRawFormat	string 		`yaml:"RAMRawFormat"`
	PulseMutedFormat string 	`yaml:"PulseMutedFormat"`
	PulseVolumeFormat string 	`yaml:"PulseVolumeFormat"`
}


/*
goroutines:
- one for reading config
- one for each module 

all the channels return to the main function where its all put together
*/
func main()  {

	var configLastModified time.Time 

    var config configInterface
	parsedConfig := config.retrieveConfig(&configLastModified)

	fmt.Println("Config file successfully loaded")
	fmt.Println(*parsedConfig)



	loopCounter := 0

	/*
	Retrieve config and get modules
	goroutine for config refresh with a loop that sleeps at whatever the refresh rate says
	goroutine all the modules where each function its own loop with a channel
	another for loop the end that loops every 1 second to update the statusbar
	*/
	modules := parsedConfig.Modules

	timeChan := make(chan string)
	// dateChan := make(chan string)
	// cpuChan := make(chan string)
	// batteryChan := make(chan string)
	// ramChan := make(chan string)
	// brightnessChan := make(chan string)
	// pulseChan := make(chan string)

	for _, module := range modules {
		switch module {
			case "time":
				go GetTime(timeChan, parsedConfig)
			case "date":

			case "mpris":

			case "cpu":

			case "battery":

			case "ram":
				
			case "brightness":

			case "pulse":
			
		}
	}

	for {
		loopCounter++
		// modules := parsedConfig.Modules

		configRefreshRate, err := strconv.Atoi(parsedConfig.RefreshConfigRate)
		if err != nil {
			configRefreshRate = 100
		}


		output := ""
		// for idx, module := range modules {
		// 	moduleData := ""
		// 	switch module {
		// 		case "time":
		// 			// moduleData += GetTime(parsedConfig.TimeFormat, 
		// 			// 						parsedConfig.TwentyFourHour)
		// 			go GetTime(timeChan, parsedConfig)
		// 			moduleData += <-timeChan
		// 		case "date":
		// 			moduleData += GetDate(parsedConfig.DateFormat)
		// 		case "mpris":
		// 			moduleData += GetMpris(parsedConfig.PlayingFormat, 
		// 									parsedConfig.PausedFormat, 
		// 									parsedConfig.MprisMaxLength, 
		// 									parsedConfig.ScrollMpris,
		// 									parsedConfig.MprisScrollSpeed)
		// 		case "cpu":
		// 			moduleData += GetCPUTemp(parsedConfig.CPUTempUnits)
		// 			// moduleData += F.GetCPUUsage()
		// 		case "battery":
		// 			moduleData += GetBatteryPercentage(parsedConfig.BatteryFormat)
		// 		case "ram":
		// 			if parsedConfig.RAMDisplay == "Percentage" {
		// 				moduleData += GetRAMUsage("format placeholder", false)
		// 			} else if parsedConfig.RAMDisplay == "Raw" {
		// 				moduleData += GetRAMData(parsedConfig.RAMRawFormat, 
		// 											parsedConfig.RAMRawUnit)
		// 			}
		// 		case "brightness":
		// 			moduleData += GetBrightness()
		// 		case "pulse":
		// 			moduleData += GetPulseVolume(parsedConfig.PulseMutedFormat, 
		// 											parsedConfig.PulseVolumeFormat)
		// 	}
		// 	if moduleData == "" {
		// 		continue
		// 	}

		// 	if idx != len(modules) - 1 {
		// 		moduleData += parsedConfig.ModuleSeperator
		// 	}

		// 	output += moduleData
		// }

		output += <-timeChan


		xsetroot := "xsetroot"
		arg1 := "-name"

		cmd := exec.Command(xsetroot, arg1, output)

		_, err = cmd.Output()

		if err != nil {
			fmt.Println(err)
			return 
		}
		if parsedConfig.RefreshConfig && loopCounter % configRefreshRate == 0 {
			parsedConfigChan := make(chan *configInterface)	
			go refreshConfig(parsedConfigChan, *parsedConfig, &configLastModified) 
			parsedConfig = <- parsedConfigChan
			fmt.Println(parsedConfig)
		}
		time.Sleep(time.Second)
	}

}



func (config *configInterface) retrieveConfig(configLastModified *time.Time) *configInterface {
	stats, err := os.Stat(fmt.Sprintf("%s/.config/go-dwm-statusbar/config.yaml", os.Getenv("HOME")))
	if err != nil {
		fmt.Println("Error with reading config file")
	}
	newModTime := stats.ModTime()
	if newModTime.IsZero() {
		*configLastModified = newModTime
	} else {
		if newModTime.After(*configLastModified) {
			*configLastModified = newModTime
		} else {
			fmt.Println("Config not modified")
			return config
		}
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/go-dwm-statusbar/config.yaml", os.Getenv("HOME")))

	fmt.Println("Updating config")
	if err != nil {
		fmt.Println("Error with reading config file")
    }

	yaml.Unmarshal(data, config)

	data = nil
	
	return config
}

func refreshConfig(config chan *configInterface, configStruct configInterface,  configLastModified *time.Time) {
		config <- configStruct.retrieveConfig(configLastModified)
}
