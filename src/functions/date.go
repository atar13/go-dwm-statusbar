package functions

import "time"

//GetDate provides the current date in the format that the user specifies 
func GetDate(format string) string {
	
	fullTime := time.Now()

	formattedDate := fullTime.Format("Jan 02")

	return formattedDate

}