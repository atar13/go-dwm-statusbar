package main

import (
	"os/exec"
	"os"
    "io/ioutil"
	"fmt"
	"strings"
	"time"
	F "./functions"
	"gopkg.in/yaml.v2"
)
//add config.json parser
//setData function 

type Time struct{
	Format string
	TwentyFourHour bool
}



type configInterface struct{
	Modules			[]string	`yaml:"Modules"`
	ModuleSeperator	string		`yaml:"ModuleSeperator"`
	TimeFormat 		string 		`yaml:"TimeFormat"`
	TwentyFourHour 	bool 		`yaml:"TwentyFourHour"`
	DateFormat		string		`yaml:"DateFormat"`
	PlayingFormat	string		`yaml:"PlayingFormat"`
	PausedFormat	string		`yaml:"PausedFormat"`

}


func main()  {

	//retrieve from config.json
	// modules := []string{"pulse","brightness", "ram", "battery", "cpu", "mpris", "time", "date"}

	desktopSession := os.Getenv("XDG_SESSION_DESKTOP")

	if(strings.Compare(desktopSession, "dwm") != 0){
		fmt.Println("Window Manager is not DWM")
		os.Exit(1)
	}

    var config configInterface
	parsedConfig := config.retrieveConfig()

	fmt.Println(parsedConfig)

	modules := parsedConfig.Modules



	for {
		output := ""
		for idx, module := range modules {
			moduleData := ""
			switch module {
				case "time":
					moduleData += F.GetTime(parsedConfig.TimeFormat, parsedConfig.TwentyFourHour)
				case "date":
					moduleData += F.GetDate(parsedConfig.DateFormat)
				case "mpris":
					moduleData += F.GetMpris(parsedConfig.PlayingFormat, parsedConfig.PausedFormat)
				case "cpu":
					moduleData += F.GetCPUTemp('F')
					moduleData += F.GetCPUUsage()
				case "battery":
					moduleData += F.GetBatteryPercentage()
				case "ram":
					// moduleData += F.GetRAMData("format placeholder", 'G')
					moduleData += F.GetRAMUsage("format placeholder", false)
				case "brightness":
					moduleData += F.GetBrightness()
				case "pulse":
					moduleData += F.GetPulseVolume()
			}
			if moduleData == "" {
				continue
			}

			if idx != len(modules) - 1 {
				//customize delimiter
				// moduleData += " | "
				moduleData += parsedConfig.ModuleSeperator
			}

			output += moduleData
		}


		xsetroot := "xsetroot"
		arg1 := "-name"

		cmd := exec.Command(xsetroot, arg1, output)

		_, err := cmd.Output()

		if err != nil {
			fmt.Println(err)
			return 
		}
		time.Sleep(100 * time.Millisecond)
	}

}



func (config *configInterface) retrieveConfig() configInterface {

    data, err := ioutil.ReadFile("config-sample.yaml")
    if err != nil {
		fmt.Println("Error with reading config file")
    }


	yaml.Unmarshal(data, config)

	return *config
}