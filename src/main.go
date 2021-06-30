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
	Modules           []string `yaml:"Modules"`
	ModuleSeperator   string   `yaml:"ModuleSeperator"`
	RefreshConfig     bool     `yaml:"RefreshConfig"`
	RefreshConfigRate string   `yaml:"RefreshConfigRate"`
	TimeFormat        string   `yaml:"TimeFormat"`
	TwentyFourHour    bool     `yaml:"TwentyFourHour"`
	DateFormat        string   `yaml:"DateFormat"`
	PlayingFormat     string   `yaml:"PlayingFormat"`
	PausedFormat      string   `yaml:"PausedFormat"`
	MprisMaxLength    string   `yaml:"MprisMaxLength"`
	ScrollMpris       bool     `yaml:"ScrollMpris"`
	MprisScrollSpeed  string   `yaml:"MprisScrollSpeed"`
	CPUTempUnits      string   `yaml:"CPUTempUnits"`
	BatteryFormat     string   `yaml:"BatteryFormat"`
	RAMDisplay        string   `yaml:"RAMDisplay"`
	RAMRawUnit        string   `yaml:"RAMRawUnit"`
	RAMRawFormat      string   `yaml:"RAMRawFormat"`
	PulseMutedFormat  string   `yaml:"PulseMutedFormat"`
	PulseVolumeFormat string   `yaml:"PulseVolumeFormat"`
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

	var tick float32 = 0

	// populates array of modules
	var modules []*moduleData
	populateModules(&modules, parsedConfig)

	// fetch new module info and initialize all goroutines
	for _, module := range modules {
		fmt.Println("making routine")
		initializeRoutine(module.name, module.channel, parsedConfig)
	}

	loopCounter++

	// loop that actually displays the data to the statusbar
	for {
		output := ""

		// loops through modules and checks if their channels have any new data
		for idx, module := range modules {
			// fmt.Println(module.name)
			select {
			case moduleOutput := <-module.channel:
				if module.name == "mpris" {
					var mprisOutput string
					maxLength := config.MprisMaxLength
					maxLengthInt, err := strconv.Atoi(maxLength)
					if err != nil {
						maxLengthInt = 1
					}
					// scrollSpeed := config.MprisScrollSpeed
					// scrollSpeedFloat, err := strconv.ParseFloat(scrollSpeed, 32)
					// if err != nil {
					// 	scrollSpeedFloat = 0.75
					// }
					if len(module.output) > maxLengthInt {
						if parsedConfig.ScrollMpris {
							tick += 0.75
							intTick := int(tick)
							startPos := intTick % len(moduleOutput)
							endPos := (intTick + maxLengthInt) % len(moduleOutput)
							// scrollSpeed := config.MprisScrollSpeed
							// scrollSpeedFloat, err := strconv.ParseFloat(scrollSpeed, 32)
							// if err != nil {
							// 	sc
							if endPos > startPos {
								mprisOutput = fmt.Sprintf("%s ", moduleOutput[startPos:endPos])
							} else {
								mprisOutput = fmt.Sprintf("%s %s", moduleOutput[startPos:], moduleOutput[:endPos])
							}

						} else {
							mprisOutput = moduleOutput[:maxLengthInt]
						}
					}
					module.output = mprisOutput 
				} else {
					module.output = moduleOutput // update module's output data if new data is received from the channel
				}
			default:
				// continue // if not then look at the next module in the array and check it's channel
			}
			if module.output != "" {
				output += module.output
				if idx != len(modules)-1 {
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
			configRefreshRate = 100
		}

		if parsedConfig.RefreshConfig && loopCounter%configRefreshRate == 0 {
			parsedConfigChan := make(chan *configInterface)
			go refreshConfig(parsedConfigChan, *parsedConfig, &configLastModified)
			parsedConfig = <-parsedConfigChan
			fmt.Println(*parsedConfig)
		}
		// time.Sleep(time.Duration(mainRefreshInterval) * time.Second)
		// maybe have half a second speed
		time.Sleep(time.Second)
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

func initializeRoutine(module string, moduleChan chan string, parsedConfig *configInterface) {
	// stops already running goroutine
	// close(moduleChan)
	// <-moduleChan
	switch module {
	case "time":
		go GetTime(moduleChan, parsedConfig)
	case "date":
		go GetDate(moduleChan, parsedConfig)
	case "mpris":
		go GetMpris(moduleChan, parsedConfig)
	case "cpu":
		go GetCPU(moduleChan, parsedConfig)
	case "battery":
		// go GetBatt
	case "ram":
		go GetRAM(moduleChan, parsedConfig)
	case "brightness":
		go GetBrightness(moduleChan, parsedConfig)
	case "pulse":
		go GetPulseVolume(moduleChan, parsedConfig)
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

func refreshConfig(config chan *configInterface, configStruct configInterface, configLastModified *time.Time) {
	config <- configStruct.retrieveConfig(configLastModified)
}
