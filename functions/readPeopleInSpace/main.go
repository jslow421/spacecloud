package main

import (
	"context"
	"encoding/json"
	"fmt"
	"functions/models"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getData() (people models.PeopleInSpace) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})

	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)

	bucketName := os.Getenv("DATA_BUCKET")
	bucketKey := os.Getenv("BUCKET_KEY")

	s3Object, s3err := downloader.S3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	})

	if s3err != nil {
		fmt.Println("Failed to retrieve file from S3", bucketName, bucketKey)
		panic(s3err)
	}

	fmt.Println("Successfully retrieved file from S3")

	defer s3Object.Body.Close()

	body, readErr := io.ReadAll(s3Object.Body)

	if readErr != nil {
		panic(readErr)
	}
	fmt.Println("Successfully read file from S3")

	fileJson := string(body)

	unmarshalErr := json.Unmarshal([]byte(fileJson), &people)

	if unmarshalErr != nil {
		panic(unmarshalErr)
	}

	fmt.Println("Successfully unmarshalled json from S3")

	return people
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	peopleInSpace := getData()

	if len(peopleInSpace.People) <= 0 {
		panic("No people in space?")
	}

	readJson, _ := json.Marshal(peopleInSpace)
	responseBody := string(readJson)

	return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       responseBody,
		},
		nil
}

func main() {
	runtime.Start(handleRequest)
}
