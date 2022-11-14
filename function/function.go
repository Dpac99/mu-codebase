package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const shimURL = "http://localhost:1234/"

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Hello!",
		StatusCode: 200,
	}, nil
}

func launch() {
	_, err := http.Post(shimURL+"register", "application/json", nil)
	if err != nil {
		log.Fatalf("Failed to register on server: %s", err)
		os.Exit(1)
	}
}

func poll() {
	stop := false
	for !stop {
		res, err := http.Get(shimURL + "poll")
		if err != nil {
			os.Exit(1)
		}
		if res != nil {
			stop = true
		}
	}
	fmt.Printf("Done")
}

func main() {
	log.Println("Hello!")
	launch()
	go poll()
	lambda.Start(HandleRequest)
}
