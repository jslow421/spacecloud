package main

import (
	"context"
	"encoding/json"
	"functions/models"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketBasics struct {
	S3Client *s3.Client
}

func getNext5Launches() models.UpcomingRocketLaunchesApiResponse {
	resp, err := http.Get("https://fdo.rocketlaunch.live/json/launches/next/5")

	if err != nil {
		log.Println("Failed to get next 5 launches")
	}

	body, bodyReadErr := io.ReadAll(resp.Body)

	if bodyReadErr != nil {
		log.Println("Failed to read response body", bodyReadErr)
	}

	var rockets = models.UpcomingRocketLaunchesApiResponse{}
	unmarshalErr := json.Unmarshal(body, &rockets)

	if unmarshalErr != nil {
		log.Println("Failed to unmarshal json", unmarshalErr)
	}

	return rockets
}

func storeNextLaunchesInS3(rockets models.UpcomingRocketLaunchesApiResponse) {
	cfg, _ := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})

	basics := BucketBasics{
		S3Client: s3.NewFromConfig(cfg),
	}

	rocketJson, jsonMarshalErr := json.Marshal(rockets)

	if jsonMarshalErr != nil {
		log.Println("Failed to marshal json", jsonMarshalErr)
	}

	reader := strings.NewReader(string(rocketJson))

	uploader := manager.NewUploader(basics.S3Client)
	_, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("DATA_BUCKET")),
		Key:    aws.String(os.Getenv("BUCKET_KEY")),
		Body:   reader,
	})

	if uploadErr != nil {
		log.Println("Failed to upload file to S3", uploadErr)
	}
}

func handler() {
	rockets := getNext5Launches()
	storeNextLaunchesInS3(rockets)
	log.Println(rockets)
}
func main() {
	runtime.Start(handler)
}
