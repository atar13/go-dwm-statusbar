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

	desktopSession := os.Getenv("XDG_SESSION_DESKTOP")

	if(strings.Compare(desktopSession, "dwm") != 0){
		fmt.Println("Window Manager is not DWM")
		os.Exit(1)
	}


	for {


		xsetroot := "xsetroot"
		arg1 := "-name"
		data := F.GetTime()

		cmd := exec.Command(xsetroot, arg1, data)

		_, err := cmd.Output()

		if err != nil {
			fmt.Println(err)
			return 
		}
		// fmt.Println("snarf")
		time.Sleep(1000 * time.Millisecond)

	}

}