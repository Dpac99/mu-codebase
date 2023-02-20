package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/struCoder/pidusage"
)

// const shimURL = "http://host.docker.internal:1234/"

const shimURL = "http://localhost:1234/"

var id string

// func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		Body:       "Hello!",
// 		StatusCode: 200,
// 	}, nil
// }

type RegisterResponse struct {
	UUID string
}

func (rr *RegisterResponse) ParseBody(body io.Reader) error {
	var uuid string
	if err := json.NewDecoder(body).Decode(&uuid); err != nil {
		return err
	}
	if len(uuid) == 0 {
		return errors.New("invalid meeting")
	}
	rr.UUID = uuid
	return nil
}

func launch() {
	var rr = &RegisterResponse{}
	res, err := http.Post(shimURL+"register", "application/json", nil)
	if err != nil {
		log.Fatalf("Failed to register on server: %s", err)
		os.Exit(1)
	}
	if err := rr.ParseBody(res.Body); err != nil {
		log.Println(res.Body)
		log.Fatalf("Failed to parse body: %s", err)
	}
	id = rr.UUID
	log.Printf("Registered with server")

}

type PollRequest struct {
	UUID   string  `json: "uuid`
	CPU    float64 `json: "cpu"`
	Memory float64 `json: "memory"`
}

func poll() {
	stop := false
	for !stop {
		log.Println("Polling...")

		sysInfo, err := pidusage.GetStat(os.Getpid())
		if err != nil {
			log.Fatalf("could not get process stat: %s", err)
		}

		pr := &PollRequest{}
		pr.CPU = sysInfo.CPU
		pr.Memory = sysInfo.Memory
		pr.UUID = id

		log.Printf("UUID: %s \nMemory: %f \nCPU:%f\n", id, sysInfo.Memory, sysInfo.CPU)
		json_data, err := json.Marshal(pr)
		if err != nil {
			log.Fatalf("could not marshall data: %s", err)
		}

		log.Println(string(json_data))

		_, err = http.Post(shimURL+"poll", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Printf("Done")
}

func main() {
	log.Println("Hello!")
	launch()
	go poll()
	// lambda.Start(HandleRequest)
	for true {
	}
}
