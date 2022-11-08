package main

import (
	"net/http"
	"os"
	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
)

const shimURL = "http://localhost:3000/register"

type MyEvent struct {
	Name string `json:"name"`
}

// func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		Body:       fmt.Sprintf("Hello!"),
// 		StatusCode: 200,
// 	}, nil
// }

func launch() {
	_, err := http.Post(shimURL, "application/json", nil)
	if err != nil {
		os.Exit(1)
	}
}
func main() {
	launch()

	// lambda.Start(HandleRequest)
}
