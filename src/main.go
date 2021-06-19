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
	// maxLoopCounter := math.Max()

	// mainRefreshInterval := 0.5

	// timeRefreshInterval := 1 / mainRefreshInterval
	// dateRefreshInterval := 10 / mainRefreshInterval
	// mprisRefreshInterval := 1 / mainRefreshInterval
	// cpuRefreshInterval := 1 / mainRefreshInterval
	// ramRefreshInterval := 1 / mainRefreshInterval

	// var displayTime string
	// var displayDate string
	// var displayMpris string
	// var displayCpu string
	// var displayRam string

	// populates array of modules
	var modules []*moduleData
	populateModules(&modules, parsedConfig)

	// getDate := math.Mod(float64(loopCounter), dateRefreshInterval) == 0
	// getTime := math.Mod(float64(loopCounter), timeRefreshInterval) == 0
	// getMpris := math.Mod(float64(loopCounter), mprisRefreshInterval) == 0
	// getCpu := math.Mod(float64(loopCounter), cpuRefreshInterval) == 0
	// getRam := math.Mod(float64(loopCounter), ramRefreshInterval) == 0

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
				module.output = moduleOutput // update module's output data if new data is received from the channel
			default:
				// continue // if not then look at the next module in the array and check it's channel
				break
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
		time.Sleep(time.Second/2)
	}

}

func populateModules(modules *[]*moduleData, parsedConfig *configInterface) {
	moduleNames := parsedConfig.Modules
	for _, moduleName := range moduleNames {
		var newModule moduleData
		newModule.name = moduleName
		mouduleChannel := make(chan string)
		newModule.channel = mouduleChannel
		// maybe check if it doesn't already exist in modules
		*modules = append(*modules, &newModule)
	}
}

func initializeRoutine(module string, moduleChan chan string, parsedConfig *configInterface) {
	// stops already running goroutine
	// close(moduleChan)
	// <-moduleChan
	switch module {
	case "time":
		// have this condition be based on time (another counter)
		// if getTime {
		fmt.Println("Getting time")
		go GetTime(moduleChan, parsedConfig)
		// only when the goroutine is called set a boolean to true meaning that the channel is ready ro recieve
		// } else {
		// if not set the boolean to false so that later it doesn't try to recieve data from the channel and instead uses previous information
		// }

	case "date":
		// if getDate {
		fmt.Println("Getting date")
		go GetDate(moduleChan, parsedConfig)
		// }
	case "mpris":
		// if getMpris {
		// 	fmt.Println("Getting Mpris")
		// 	go GetMpris(moduleChan, parsedConfig)
		// }
	case "cpu":
		// if getCpu {
		// 	fmt.Println("Getting Cpu")
		// 	go GetCPU(moduleChan, parsedConfig)
		// }
	case "battery":

	case "ram":
		// if getRam {
		// 	fmt.Println("Getting RAM")
		// 	go GetRAM(moduleChan, parsedConfig)
		// }
	case "brightness":

	case "pulse":

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
