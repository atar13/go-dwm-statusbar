package main

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

//GetCPUTemp returns the current CPU temperature as specified in the units passed as a parameter
func GetCPUTemp(config *configInterface) string {
	unit := config.CPUTempUnits

	cat := "cat"
	arg1 := "/sys/class/thermal/thermal_zone0/temp"

	cmd := exec.Command(cat, arg1)

	temp, err := cmd.Output()

	if err != nil {
		return ""
	}
	tempString := string(temp[:len(temp)-1])
	tempInt, err := strconv.ParseFloat(tempString, 10)
	if err != nil {
		return ""
	}

	tempInt /= 1000

	if unit == "C" {
		return fmt.Sprintf("%.1f°C", math.Ceil(tempInt))
	} else if unit == "F" {
		tempInt = (tempInt * 1.8) + 32
		return fmt.Sprintf("%.1f°F", math.Ceil(tempInt))
	} else {
		return ""
	}
}

/*
user    nice   system  idle      iowait irq   softirq  steal  guest  guest_nice
cpu  74608   2520   24433   1117073   6176   4054  0        0      0      0

PrevIdle = previdle + previowait
Idle = idle + iowait

PrevNonIdle = prevuser + prevnice + prevsystem + previrq + prevsoftirq + prevsteal
NonIdle = user + nice + system + irq + softirq + steal

PrevTotal = PrevIdle + PrevNonIdle
Total = Idle + NonIdle

# differentiate: actual value minus the previous one
totald = Total - PrevTotal
idled = Idle - PrevIdle

CPU_Percentage = (totald - idled)/totald
*/

//
func GetCPUUsage() string {

	// var previdle int
	// var previoWait int
	// var PrevIdle int

	// var idle int
	// var iowait int
	// var Idle int

	// var prevuser int
	// var prevnice int64
	// var prevsystem int64
	// var previrq int
	// var prevsoftirq int
	// var prevsteal int
	// var PrevNonIdle int

	// var user int
	// var nice int
	// var system int
	// var irq int
	// var softirq int
	// var steal int
	// var NonIdle int

	// var PrevTotal int

	// var Total int

	// var totald int

	// var idled int

	// var percentage int

	// for j := 0; j < 2; j++ {

	// 	cat := "cat"
	// 	arg1 := "/proc/stat"

	// 	cmd := exec.Command(cat, arg1)

	// 	byteData, err := cmd.Output()
	// 	if err != nil {
	// 		return ""
	// 	}

	// 	data := string(byteData)

	// 	var s scanner.Scanner
	// 	s.Init(strings.NewReader(data))
	// 	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
	// 		if(s.TokenText()=="cpu"){
	// 			fmt.Println(s.Position, s.TokenText())
	// 			for i := 0; i < 10; i++ {
	// 				s.Scan()
	// 				nextWord := s.TokenText()
	// 				switch j {
	// 				case 0:
	// 					switch i {
	// 					case 0:
	// 						prevnice, err = strconv.ParseInt(nextWord, 10, 32)
	// 						if err != nil {
	// 							return ""
	// 						}
	// 					case 1:
	// 						prevsystem, err = strconv.ParseInt(nextWord, 10, 32)
	// 						if err != nil {
	// 							return ""
	// 						}

	// 					}
	// 				case 1:

	// 				}
	// 				if s.TokenText()=="cpu0" {
	// 					break
	// 				}
	// 			}
	// 			break
	// 		}
	// 	}

	// 	if j == 0 {
	// 		//pause for half a second
	// 	}

	// }

	// // for tok :=

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("%o", int(math.Ceil(percent[0])))
}

func GetCPU(cpuChan chan string, config *configInterface) {
	defaultFormat := "Temp:%s Usage:%s%%"
	cpuChan <- fmt.Sprintf(defaultFormat, GetCPUTemp(config), GetCPUUsage())
}
