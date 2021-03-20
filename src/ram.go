package main

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
func GetRAMData(format string, unit string) string {
	updateRAMStatus()

	var usedString string
	var freeString string
	var availableString string
	var totalString string

	if unit == "G" {
		availableGig := float64(available)/1000
		usedGig := float64(used)/1000
		freeGig := float64(free)/1000
		totalGig := float64(total)/1000

		availableString = fmt.Sprintf("%0.02f", availableGig)
		usedString = fmt.Sprintf("%0.02f", usedGig)
		freeString = fmt.Sprintf("%0.02f", freeGig)
		totalString = fmt.Sprintf("%0.02f", totalGig)

	} else {
		availableString = fmt.Sprintf("%v", int64(available))
		usedString = fmt.Sprintf("%v", int64(used))
		freeString = fmt.Sprintf("%v", int64(free))
		totalString = fmt.Sprintf("%v", int64(total))
	}

	output := ""

	for i := 0; i < len(format); i++ {
		char := string(format[i])
		
		if char == "@" {
			i++;
			nextChar := string(format[i])

			switch nextChar {
			case "u":
				output += usedString;
			case "a":
				output += availableString;
			case "f":
				output += freeString;
			case "t":
				output += totalString;
			}
		} else {
			output += char
		}
	}


	return output
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
