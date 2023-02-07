package main

import (
	"encoding/json"
	"fmt"
	"functions/models"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
)

func getNext5Launches() models.UpcomingRocketLaunchesApiResponse {
	resp, err := http.Get("https://fdo.rocketlaunch.live/json/launches/next/5")

	if err != nil {
		fmt.Println("Failed to get next 5 launches")
	}

	body, bodyReadErr := io.ReadAll(resp.Body)

	if bodyReadErr != nil {
		fmt.Println("Failed to read response body", bodyReadErr)
	}

	var rockets = models.UpcomingRocketLaunchesApiResponse{}
	unmarshalErr := json.Unmarshal(body, &rockets)

	if unmarshalErr != nil {
		fmt.Println("Failed to unmarshal json", unmarshalErr)
	}

	return rockets
}
func handler() {
	rockets := getNext5Launches()
	fmt.Println(rockets)
}
func main() {
	runtime.Start(handler)
}
