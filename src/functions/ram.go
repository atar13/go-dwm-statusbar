package functions

import (
	"os/exec"
	"strings"
	"text/scanner"
	"strconv"
	"fmt"
)

var available int64
var used int64
var total int64
var free int64

//format can use available/used/total/free
func updateRAMStatus() {
	app := "free"
	arg1 := "-m"
	cmd := exec.Command(app, arg1)
	output, err := cmd.Output()
	if err != nil {
		available = 0
		used = 0
		total = 0 
		free = 0
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(string(output)))
	for i := 0; i < 16; i++ {
		s.Scan()
		switch i {
		case 10:
			totalString := s.TokenText()
			total, err = strconv.ParseInt(totalString, 10, 32)
			if err != nil {
				total = 0
			}
		case 11:
			usedString := s.TokenText()
			used, err = strconv.ParseInt(usedString, 10, 32)
			if err != nil {
				used = 0
			}
		case 12:
			freeString := s.TokenText()
			free, err = strconv.ParseInt(freeString, 10, 32)
			if err != nil {
				free = 0
			}
		case 15:
			availableString := s.TokenText()
			available, err = strconv.ParseInt(availableString, 10, 32)
			if err != nil {
				available = 0
			}
		}
	}


}

//GetRAMData provides the ...
func GetRAMData(format string, unit rune) string {
	updateRAMStatus()
	if unit == 'G' {
		availableGig := float64(available)/1000
		usedGig := float64(used)/1000
		freeGig := float64(free)/1000
		totalGig := float64(total)/1000

		return fmt.Sprintf("%0.02f %0.02f %0.02f %0.02f",availableGig, usedGig, freeGig, totalGig)
	}

	placeholder := fmt.Sprint(free, available, used, total)
	return placeholder
}

//take care of divide by zero error
func GetRAMUsage(format string, useFree bool) string {
	updateRAMStatus()
	var percentage float64
	if useFree {
		percentage = float64(free)/float64(available)
	} else {
		percentage = float64(used)/float64(available)
	}
	percentage *= 100
	return fmt.Sprintf("ram: %0.01f%%", percentage)
}
