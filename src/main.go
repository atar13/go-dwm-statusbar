package main

import "os/exec"
import "os"
import "fmt"
import "strings"

import "time"
import (
    F "./functions"
)
//add config.json parser
//setData function 




func main()  {
	
	//retreive from config.json
	modules := []string{"time", "date"}

	desktopSession := os.Getenv("XDG_SESSION_DESKTOP")

	if(strings.Compare(desktopSession, "dwm") != 0){
		fmt.Println("Window Manager is not DWM")
		os.Exit(1)
	}


	for {

		output := ""

		for idx, module := range modules {
			switch module {
				case "time":
					output += F.GetTime()
				case "date":
					output += F.GetDate("format placeholder")
			}

			if(idx != len(modules) - 1){
				output += " | "
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
		time.Sleep(1000 * time.Millisecond)
	}

}