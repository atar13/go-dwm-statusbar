package main

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"time"
	F "./functions"
)
//add config.json parser
//setData function 




func main()  {
	
	//retrieve from config.json
	modules := []string{"ram", "battery", "cpu", "mpris", "time", "date"}

	desktopSession := os.Getenv("XDG_SESSION_DESKTOP")

	if(strings.Compare(desktopSession, "dwm") != 0){
		fmt.Println("Window Manager is not DWM")
		os.Exit(1)
	}


	for {
		output := ""
		for idx, module := range modules {
			moduleData := ""
			switch module {
				case "time":
					moduleData += F.GetTime()
				case "date":
					moduleData += F.GetDate("format placeholder")
				case "mpris":
					moduleData += F.GetMpris()
				case "cpu":
					moduleData += F.GetCPUTemp('F')
					moduleData += F.GetCPUUsage()
				case "battery":
					moduleData += F.GetBatteryPercentage()
				case "ram":
					moduleData += F.GetRAMData("format placeholder", 'G')
					// moduleData += F.GetRAMUsage("format placeholder", false)

			}
			if moduleData == "" {
				continue
			}

			if idx != len(modules) - 1 {
				//customize delimiter
				moduleData += " | "
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
		time.Sleep(1000 * time.Millisecond)
	}

}