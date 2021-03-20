package main 

import (
	"net/http"
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

func GetCurrentAQI(location string) string {

	err := godotenv.Load();
	if err != nil {
		fmt.Println("Error loading .env file");
	}

	token := os.Getenv("AQI_TOKEN")
	apiURL := fmt.Sprintf("https://aqi.waqi.info/feed/%v/?token=%v", location, token);
	// apiURL := "https://api.waqi.info/feed/beijing/?token=1028a7ca49148c0640560bcc20ce877b38c4e37a"
	fmt.Println(apiURL);
	res, err := http.Get(apiURL);

	if err != nil {
		fmt.Println("Error with retrieving AQI Info");
	}

	fmt.Println(res);

	return ""
}