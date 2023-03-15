package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"serverless/tasks"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/struCoder/pidusage"
)

const shimURL = "http://host.docker.internal:1234/"

var id string

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
	for {

		sysInfo, err := pidusage.GetStat(os.Getpid())
		if err != nil {
			log.Fatalf("could not get process stat: %s", err)
		}

		pr := &PollRequest{}
		pr.CPU = sysInfo.CPU
		pr.Memory = sysInfo.Memory
		pr.UUID = id

		log.Printf("UUID: %s \nMemory: %f \nCPU:%f\n", pr.UUID, pr.Memory, pr.CPU)
		json_data, err := json.Marshal(pr)
		if err != nil {
			log.Fatalf("could not marshall data: %s", err)
		}

		_, err = http.Post(shimURL+"poll", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			log.Println("err")
			os.Exit(1)
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {
	log.Println("Hello!")
	launch()
	go poll()
	lambda.Start(tasks.ExecuteTask)
}
