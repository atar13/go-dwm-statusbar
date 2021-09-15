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

type configInterface struct {
	Modules              []string `yaml:"Modules"`
	ModuleSeperator      string   `yaml:"ModuleSeperator"`
	RefreshConfig        bool     `yaml:"RefreshConfig"`
	RefreshConfigRate    string   `yaml:"RefreshConfigRate"`
	TimeFormat           string   `yaml:"TimeFormat"`
	TwentyFourHour       bool     `yaml:"TwentyFourHour"`
	DateFormat           string   `yaml:"DateFormat"`
	PlayingFormat        string   `yaml:"PlayingFormat"`
	PausedFormat         string   `yaml:"PausedFormat"`
	MprisMaxLength       string   `yaml:"MprisMaxLength"`
	ScrollMpris          bool     `yaml:"ScrollMpris"`
	MprisScrollSpeed     string   `yaml:"MprisScrollSpeed"`
	CPUTempUnits         string   `yaml:"CPUTempUnits"`
	BatteryFormat        string   `yaml:"BatteryFormat"`
	ChargingIndicator    string   `yaml:"ChargingIndicator"`
	DischargingIndicator string   `yaml:"DischargingIndicator"`
	RAMDisplay           string   `yaml:"RAMDisplay"`
	RAMRawUnit           string   `yaml:"RAMRawUnit"`
	RAMRawFormat         string   `yaml:"RAMRawFormat"`
	LowBrightnessFormat  string   `yaml:"LowBrightnessFormat"`
	HighBrightnessFormat string   `yaml:"HighBrightnessFormat"`
	PulseMutedFormat     string   `yaml:"PulseMutedFormat"`
	PulseVolumeFormat    string   `yaml:"PulseVolumeFormat"`
}

type moduleData struct {
	name    string
	channel chan string
	output  string
}

/*
goroutines:
- one for reading config
- one for each module

all the channels return to the main function where its all put together
*/
func main() {

	var configLastModified time.Time

	var config configInterface
	var parsedConfig *configInterface
	parsedConfig = config.retrieveConfig(&configLastModified)

	fmt.Println("Config file successfully loaded")
	fmt.Println(*parsedConfig)

	loopCounter := 0

	// var tick float32 = 0

	// populates array of modules
	// var modules []*moduleData
	// populateModules(&modules, parsedConfig)

	// // fetch new module info and initialize all goroutines
	// for _, module := range modules {
	// 	fmt.Printf("making %v routine\n", module.name)
	// 	initializeRoutine(module.name, module.channel, parsedConfig)
	// }

	// loopCounter++

	tick := 0

	// timeOutput := ""
	dateOutput := ""
	cpuOutput := ""
	ramOutput := ""
	batteryOutput := ""

	// loop that actually displays the data to the statusbar
	for {
		output := ""

		// loops through modules and checks if their channels have any new data
		for idx, module := range parsedConfig.Modules {
			moduleOutput := ""
			switch module {
			case "time":
				moduleOutput = GetTime(parsedConfig)
			case "date":
				if loopCounter%20 == 0 {
					moduleOutput = GetDate(parsedConfig)
					dateOutput = moduleOutput
				} else {
					moduleOutput = dateOutput
				}
			case "cpu":
				if loopCounter%20 == 0 {
					moduleOutput = GetCPUUsage(parsedConfig)
					cpuOutput = moduleOutput
				} else {
					moduleOutput = cpuOutput
				}
			case "ram":
				if loopCounter%4 == 0 {
					moduleOutput = GetRAMUsage(parsedConfig)
					ramOutput = moduleOutput
				} else {
					moduleOutput = ramOutput
				}
			case "battery":
				if loopCounter%4 == 0 {
					moduleOutput = GetBatteryPercentage(parsedConfig)
					batteryOutput = moduleOutput
				} else {
					moduleOutput = batteryOutput
				}
			case "brightness":
				moduleOutput = GetBrightness(parsedConfig)
			case "pulse":
				moduleOutput = GetPulseVolume(parsedConfig)
			case "mpris":
				moduleOutput = GetMpris(parsedConfig)
				maxLength, err := strconv.Atoi(config.MprisMaxLength)
				if err != nil {
					maxLength = 10
				}
				if len(moduleOutput) > maxLength {
					// each iteration move the text to 1/2 of the maxlength
					tick += maxLength / 2
					startPos := tick % len(moduleOutput)
					endPos := (tick + maxLength) % len(moduleOutput)
					if endPos > startPos {
						moduleOutput = fmt.Sprintf("%s ", moduleOutput[startPos:endPos])
					} else {
						moduleOutput = fmt.Sprintf("%s %s", moduleOutput[startPos:], moduleOutput[:endPos])
					}

				}
			}
			if moduleOutput != "" {
				output += moduleOutput
				if idx != len(parsedConfig.Modules)-1 {
					output += parsedConfig.ModuleSeperator
				}
			}
		}

		xsetroot := "xsetroot"
		arg1 := "-name"

		cmd := exec.Command(xsetroot, arg1, output)

		_, err := cmd.Output()

		if err != nil {
			fmt.Println(err)
			return
		}

		configRefreshRate, err := strconv.Atoi(parsedConfig.RefreshConfigRate)
		if err != nil {
			configRefreshRate = 10
		}

		if parsedConfig.RefreshConfig && loopCounter%configRefreshRate == 0 {
			parsedConfigChan := make(chan *configInterface)
			go refreshConfig(parsedConfigChan, *parsedConfig, &configLastModified)
			parsedConfig = <-parsedConfigChan
			fmt.Println("Reading config file")
			fmt.Println(*parsedConfig)
		}
		// time.Sleep(time.Duration(mainRefreshInterval) * time.Second)
		// maybe have half a second speed
		time.Sleep(1 * time.Second)
		loopCounter++
	}

}

func populateModules(modules *[]*moduleData, parsedConfig *configInterface) {
	moduleNames := parsedConfig.Modules
	for _, moduleName := range moduleNames {
		var newModule moduleData
		newModule.name = moduleName
		mouduleChannel := make(chan string)
		newModule.channel = mouduleChannel
		// TODO: check if it doesn't already exist in modules
		*modules = append(*modules, &newModule)
	}
}

// func initializeRoutine(module string, moduleChan chan string, parsedConfig *configInterface) {
// 	// stops already running goroutine
// 	switch module {
// 	case "time":
// 		go GetTime(moduleChan, parsedConfig)
// 	case "date":
// 		go GetDate(moduleChan, parsedConfig)
// 	case "mpris":
// 		go GetMpris(moduleChan, parsedConfig)
// 	case "cpu":
// 		go GetCPU(moduleChan, parsedConfig)
// 	case "battery":
// 		go GetBatteryPercentage(moduleChan, parsedConfig)
// 	case "ram":
// 		go GetRAM(moduleChan, parsedConfig)
// 	case "brightness":
// 		go GetBrightness(moduleChan, parsedConfig)
// 	case "pulse":
// 		go GetPulseVolume(moduleChan, parsedConfig)
// 	}
// }

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

	// data = nil

	return config
}

func refreshConfig(config chan *configInterface, configStruct configInterface, configLastModified *time.Time) {
	config <- configStruct.retrieveConfig(configLastModified)
}
