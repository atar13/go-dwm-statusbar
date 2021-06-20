package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

var tick float32 = 0

/*
GetMpris returns the ...
Compatible with cmus, vlc, and partially spotify
POSITION DOESNT WORK ON SPOTIFY, IT ALWAYS DISPLAYS ZERO
TODO: custom formatting for play state and pause state with a parser that converts it to a string ready to pass into fmt.sprintf
*/
func GetMpris(mprisChan chan string, config *configInterface) {
	// func GetMpris(playingFormat string, pausedFormat string, maxLength string, scroll bool, scrollSpeed string) string {
	playingFormat := config.PlayingFormat
	pausedFormat := config.PausedFormat
	maxLength := config.MprisMaxLength
	scroll := config.ScrollMpris
	scrollSpeed := config.MprisScrollSpeed

	interval := time.Second

	for {
		con, conErr := dbus.SessionBus()

		if conErr != nil {
			fmt.Println("Error with connecting to DBUS")
			mprisChan <- ""
			time.Sleep(interval)
			continue
		}
		players, playerErr := mpris.List(con)

		if playerErr != nil {
			fmt.Println("Error with detecting player")
			mprisChan <- ""
			time.Sleep(interval)
			continue
		}
		if len(players) == 0 {
			// fmt.Println("No player found")
			mprisChan <- ""
			time.Sleep(interval)
			continue
		}

		name := players[0]
		player := mpris.New(con, name)

		status, err := player.GetPlaybackStatus()
		if err != nil {
			mprisChan <- ""
			time.Sleep(interval)
			continue
		}


		scrollSpeedFloat, err := strconv.ParseFloat(scrollSpeed, 32)
		if err != nil {
			scrollSpeedFloat = 0.75
		}

		if status == "Playing" {
			mprisChan <- getPlayingInfo(player, playingFormat, maxLengthInt, scroll, float32(scrollSpeedFloat))
			time.Sleep(interval)
			continue
		} else if status == "Paused" {
			//have an option to format pause state
			quit := false
			for _, player := range players {
				player := mpris.New(con, player)

				status, err := player.GetPlaybackStatus()
				if err != nil {
					continue
				}
				if status == "Playing" {
					mprisChan <- getPlayingInfo(player, playingFormat, maxLengthInt, scroll, float32(scrollSpeedFloat))
					quit = true
					break
				}
			}
			if quit {
				time.Sleep(interval)
				continue
			}
			mprisChan <- pausedFormat
			time.Sleep(interval)
			continue
		} else {
			mprisChan <- ""
			time.Sleep(interval)
			continue

		}
	}
}

func getPlayingInfo(player *mpris.Player, playingFormat string, maxLength int, scroll bool, scrollSpeed float32) string {
	metadata, err := player.GetMetadata()
	if err != nil {
		return ""
	}

	title := ""
	album := ""
	artist := ""
	albumArtist := ""
	formattedTrackLength := ""
	formattedPosition := ""

	for key := range metadata {
		switch key {
		case "xesam:title":
			title = metadata["xesam:title"].String()
			title = title[1 : len(title)-1]
		case "xesam:album":
			album = metadata["xesam:album"].String()
			album = album[1 : len(album)-1]
		case "xesam:artist":
			artist = metadata["xesam:artist"].String()
			artist = artist[2 : len(artist)-2]
		case "xesam:albumArtist":
			albumArtist = metadata["xesam:albumArtist"].String()
			albumArtist = albumArtist[2 : len(albumArtist)-2]
		case "mpris:length":
			trackLength := metadata["mpris:length"].String()
			if len(trackLength) == 0 {
				return ""
			}
			trackLengthInt, err := strconv.ParseInt(trackLength[3:], 10, 64)
			trackLengthInt /= 1000000
			if err != nil {
				return ""
			}
			trackLengthMin := trackLengthInt / 60
			trackLengthSec := trackLengthInt % 60
			formattedTrackLength = fmt.Sprintf("%v:%02d", trackLengthMin, trackLengthSec)
		}
	}

	position, err := player.GetPosition()
	if err != nil {
		return ""
	}
	positionInt := int(position)
	positionMin := positionInt / 60
	positionSec := positionInt % 60
	formattedPosition = fmt.Sprintf("%v:%02d", positionMin, positionSec)

	// fmt.Println(album, artist, formattedPosition, formattedTrackLength)

	output := ""

	if playingFormat == "" {
		output = title + " by " + albumArtist
		return "▶ Playing: " + output
	}

	// TODO: change this to strings.replace
	for i := 0; i < len(playingFormat); i++ {
		char := string(playingFormat[i])
		if char == "@" {
			nextChar := string(playingFormat[i+1])
			// fmt.Println(nextChar)
			i++
			switch nextChar {
			case "t":
				output += title
			case "p":
				output += formattedPosition
			case "l":
				output += formattedTrackLength
			case "a":
				// skipNextChar = true
				if string(playingFormat[i+1]) == "r" {
					if string(playingFormat[i+2]) == "t" {
						output += artist
						i += 2
					}
				} else if string(playingFormat[i+1]) == "l" {
					if string(playingFormat[i+2]) == "b" {
						output += album
						i += 2
					} else if string(playingFormat[i+2]) == "a" {
						if string(playingFormat[i+3]) == "r" {
							output += albumArtist
							i += 3
						}
					}
				} else {

				}
			}
		} else {
			output += string(char)
		}
	}
	fmt.Println(output)

	return output
}
