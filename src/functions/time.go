package functions

import "time"
// import "fmt"

//GetTime retrieves the current time and formats it according to the user's configuration
func GetTime() string{

	fullTime := time.Now()
	// data := time.Format() time.Now().String()

	// data := fmt.Sprintf("%v:%v:%v", currentTime.Hour(),currentTime.Minute(), currentTime.Second())
	formattedTime := fullTime.Format("03:04:05")
	return formattedTime	
}

