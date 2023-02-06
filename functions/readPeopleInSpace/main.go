package main

import (
	"context"
	"encoding/json"
	"functions/models"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PersonApiResponse struct {
	Message string          `json:"message"`
	People  []models.Person `json:"people"`
}

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

	s3Object, s3err := downloader.S3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("people-in-space"),
		Key:    aws.String("people-in-space.json"),
	})

	if s3err != nil {
		panic(s3err)
	}

	return
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	peopleInSpace := getData()
	value := ""

	if len(peopleInSpace.People) > 0 {
		readJson, _ := json.Marshal(peopleInSpace)
		value = string(readJson)
	}

	return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       value,
		},
		nil
}

func main() {
	runtime.Start(handleRequest)
}
