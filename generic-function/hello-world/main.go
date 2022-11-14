// package main

// import (
// 	"errors"
//

// 	"io/ioutil"
// 	"net/http"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// )

// var (
// 	// DefaultHTTPGetAddress Default Address
// 	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

// 	// ErrNoIP No IP found in response
// 	ErrNoIP = errors.New("No IP in HTTP response")

// 	// ErrNon200Response non 200 status code in response
// 	ErrNon200Response = errors.New("Non 200 Response found")
// )

// func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	resp, err := http.Get(DefaultHTTPGetAddress)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{}, err
// 	}

// 	if resp.StatusCode != 200 {
// 		return events.APIGatewayProxyResponse{}, ErrNon200Response
// 	}

// 	ip, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{}, err
// 	}

// 	if len(ip) == 0 {
// 		return events.APIGatewayProxyResponse{}, ErrNoIP
// 	}

// 	return events.APIGatewayProxyResponse{
// 		Body:       fmt.Sprintf("Hello, %v", string(ip)),
// 		StatusCode: 200,
// 	}, nil
// }

// func main() {
// 	lambda.Start(handler)
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const shimURL = "http://192.168.1.90g:1234/"

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello!"),
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
	launch()
	go poll()
	lambda.Start(HandleRequest)
}
