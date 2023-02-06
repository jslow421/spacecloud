package main

import (
	"context"
	"encoding/json"
	"fmt"
	"functions/models"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketBasics struct {
	S3Client *s3.Client
}

func getPeopleInSpaceFromApi() (people models.PersonInSpaceApiResponse) {
	resp, err := http.Get("http://api.open-notify.org/astros.json")

	if err == nil {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Failed to read response body ", err)
			defer resp.Body.Close()
		}

		unMarshalErr := json.Unmarshal(body, &people)

		if unMarshalErr != nil {
			fmt.Println("Failed to unmarshal json ", unMarshalErr)
		}
	}

	return people
}

func (basics BucketBasics) writeJsonToS3(people []models.Person) error {

	// Create update model
	update := models.PersonInSpace{
		UpdatedTime: time.Now().String(),
		People:      people,
	}

	jsonResponse, jsonMarshalErr := json.Marshal(update)
	if jsonMarshalErr != nil {
		return jsonMarshalErr
	}

	fmt.Println("JSON value: ", string(jsonResponse))

	reader := strings.NewReader(string(jsonResponse))

	uploader := manager.NewUploader(basics.S3Client)
	_, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("DATA_BUCKET")),
		Key:    aws.String(os.Getenv("BUCKET_KEY")),
		Body:   reader,
	})

	if uploadErr != nil {
		log.Println("Failed to upload file to S3", uploadErr)
	}

	return nil
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})

	if err != nil {
		panic(err)
	}

	basics := BucketBasics{
		S3Client: s3.NewFromConfig(cfg),
	}

	// Retrieve people in space from API
	people := getPeopleInSpaceFromApi()

	// write text file to S3
	writeErr := basics.writeJsonToS3(people.People)

	if writeErr != nil {
		fmt.Println("Failed to write people in space to S3", err)
		panic(err)
	}

	fmt.Println("Successfully wrote people in space to S3")
	return true, nil
}

func main() {
	fmt.Println("Starting people in space lambda function")
	runtime.Start(handleRequest)
}
