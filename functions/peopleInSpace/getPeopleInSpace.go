package peopleInSpace

import (
	"context"
	"encoding/json"
	"functions/models"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
)

type PersonApiResponse struct {
	Message string          `json:"message"`
	People  []models.Person `json:"people"`
}

func getData() (people []models.Person) {
	resp, err := http.Get("http://api.open-notify.org/astros.json")

	if err == nil {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			defer resp.Body.Close()
			return
		}

		unMarshalErr := json.Unmarshal(body, &people)

		if unMarshalErr != nil {
			return []models.Person{}
		}
	}

	return
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "",
		},
		nil
}

func main() {
	runtime.Start(handleRequest)
}
