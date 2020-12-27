package functions

import (
	"fmt"
	"strconv"
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus"
)

/*
GetMpris returns the ... 
Compatible with cmus, vlc, and partially spotify
POSITION DOESNT WORK ON SPOTIFY, IT ALWAYS DISPLAYS ZERO
TODO: custom formatting for play state and pause state with a parser that converts it to a string ready to pass into fmt.sprintf

*/
func GetMpris() string {
	con, conErr := dbus.SessionBus()

	if conErr != nil {
		fmt.Println("Error with connecting to DBUS")
		return ""
	}
	players, playerErr := mpris.List(con)

	if playerErr != nil {
		fmt.Println("Error with detecting player")
		return ""
	}
	if len(players) == 0 {
		fmt.Println("No player found")
		return ""
	}

	name := players[0]
	player := mpris.New(con, name)

	status, err := player.GetPlaybackStatus()
	if err != nil {
		return ""
	}

	if status == "Playing" {
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
				title = title[1:len(title)-1]
			case "xesam:album":
				album = metadata["xesam:album"].String()
				album = album[1:len(album)-1]
			case "xesam:artist":
				artist = metadata["xesam:artist"].String()
				artist = artist[2:len(artist)-2]
			case "xesam:albumArtist":
				albumArtist = metadata["xesam:albumArtist"].String()
				albumArtist =  albumArtist[2:len(albumArtist)-2]
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
				trackLengthMin := trackLengthInt/60
				trackLengthSec := trackLengthInt%60
				formattedTrackLength = fmt.Sprintf("%v:%02d",trackLengthMin, trackLengthSec)
			}
		}

		position, err := player.GetPosition()
		if err != nil {
			return ""
		}
		positionInt  := int(position)
		positionMin := positionInt/60
		positionSec := positionInt%60
		formattedPosition = fmt.Sprintf("%v:%02d",positionMin, positionSec)

		fmt.Println(album, artist, formattedPosition, formattedTrackLength)


		output := title + " by " + albumArtist
		return "▶ Playing: " + output

	} else if status == "Paused" {
		//have an option to format pause state
		return "▌▌ Paused"
	} else {
		return ""
	}
	
	

}