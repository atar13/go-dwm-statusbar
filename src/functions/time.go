package functions

import "time"
import "fmt"

func GetTime() string{

	currentTime := time.Now()
	// data := time.Format() time.Now().String()

	data := fmt.Sprintf("%v:%v:%v", currentTime.Hour(),currentTime.Minute(), currentTime.Second())
	return data
}

