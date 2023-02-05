package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"functions/models"
	"io"
	"log"
	"net/http"
	"os"
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
	// Set up client
	//cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	//client := s3.NewFromConfig(cfg)

	// Create update model
	update := models.PersonInSpace{
		UpdatedTime: time.Now().String(),
		People:      people,
	}

	jsonResponse, jsonMarshalErr := json.Marshal(update)
	if jsonMarshalErr != nil {
		return jsonMarshalErr
	}

	// write jsonResponse to temp file
	fileName := "peopleInSpace.json"
	file, fileCreateErr := os.Create("/tmp/" + fileName)
	if fileCreateErr != nil {
		fmt.Println("Failed to create file ", fileCreateErr)
	}

	fmt.Println("JSON value: ", string(jsonResponse))

	w := bufio.NewWriter(file)
	bytesWritten, writeErr := w.Write(jsonResponse)
	w.Flush()
	//bytesWritten, writeErr := file.Write(jsonResponse)

	if writeErr != nil {
		fmt.Println("Failed to write file ")
		panic(writeErr)
	}

	fmt.Println("Wrote ", bytesWritten, " bytes to file")

	defer file.Close()
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		fmt.Println("Failed to close file ", err)
	//		panic(err)
	//	}
	//	fmt.Println("Closed file")
	//}(file)

	// Instantiate uploader
	//uploader := manager.NewUploader(client)

	// Upload created json file to s3
	//_, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
	//	Bucket: aws.String(os.Getenv("DATA_BUCKET")),
	//	Key:    aws.String(os.Getenv("BUCKET_KEY")),
	//	Body:   file,
	//})

	_, putErr := basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("DATA_BUCKET")),
		Key:         aws.String(os.Getenv("BUCKET_KEY")),
		ContentType: aws.String("application/json"),
		Body:        file,
	})

	if putErr != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			fileName, os.Getenv("DATA_BUCKET"), os.Getenv("BUCKET_KEY"), jsonMarshalErr)
	}

	//upParams := &s3manager.UploadInput{
	//	Bucket: &bucketName,
	//	Key:    &keyName,
	//	Body:   file,
	//}

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
