package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"serverless/tasks"
	"serverless/types"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pbnjay/memory"
	"github.com/struCoder/pidusage"
)

var shimURL = "http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com"

var totalMem = float64(memory.TotalMemory())

var id string

type T = struct{}

var end_channel = make(chan T)

func Listen(req events.LambdaFunctionURLRequest) (string, error) {
	launch()
	go poll()
	end_channel <- struct{}{}
	return id, nil
}

func launch() {
	var rr = &types.RegisterResponse{}
	res, err := http.Post(shimURL+"/register", "application/json", nil)
	if err != nil {
		log.Fatalf("Could not communicate with coordinator")
	}
	res.Close = true
	defer res.Body.Close()

	if err != nil {
		log.Fatalf("Failed to register on server: %s\n", err)
		os.Exit(1)
	}
	if err := rr.ParseBody(res.Body); err != nil {
		log.Println(res.Body)
		log.Fatalf("Failed to parse body: %s\n", err)
	}
	id = rr.UUID

}

/* Sends resource usage data to coordinator as well as calculated return values */
func poll() {
	var wg sync.WaitGroup
	end := false
	for !end {

		log.Println("Polling")

		/* START collecting data */
		sysInfo, err := pidusage.GetStat(os.Getpid())
		if err != nil {
			log.Fatalf("could not get process stat: %s\n", err)
		}
		pr := &types.PollRequest{}
		pr.CPU = sysInfo.CPU
		pr.Memory = (sysInfo.Memory / totalMem) * 100
		pr.UUID = id

		json_data, err := json.Marshal(pr)
		if err != nil {
			log.Fatalf("could not marshall data: %s\n", err)
		}
		/* END collecting data */

		/* Send data to coordinator */
		res, err := http.Post(shimURL+"/poll", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatalf("Error sending data to central: %s\n", err)
			res.Body.Close()
		}

		pollResponse := types.TaskRequest{}

		if err := pollResponse.ParseBody(res.Body); err != nil {
			log.Fatalf("Error parsing poll response: %s\n", err)
			res.Body.Close()
		}

		/* If no request to execute, end */
		if pollResponse.ID == "-1" {
			log.Println("Received shutdown signal")
			end = true
			wg.Wait()
			log.Println("Ending poll")
			<-end_channel
			res.Body.Close()
		} else if pollResponse.ID == "0" {
			log.Println("Continuing as normal")
			// Continue polling
		} else {
			log.Println("Received new Task")
			wg.Add(1) // Increment WaitGroup

			/* Send request to executor, defer WaitGroup finish*/
			go func() {
				defer wg.Done()
				log.Println("Executing task")
				data, err := tasks.ExecuteTask(pollResponse)
				log.Println("Task done executing")
				res := &types.TaskResult{}
				res.ID = id
				log.Println(id)

				if err != nil {
					e, err2 := json.Marshal(err)
					if err2 != nil {
						log.Fatalf("Error marshaling data: %s\n", err2)
					}
					res.Data = e
				} else {
					res.Data = data
				}
				json_data, err := json.Marshal(res)

				if err != nil {
					log.Fatalf("could not marshall data: %s\n", err)
				}
				_, err = http.Post(shimURL+"/sendResult/"+pollResponse.ID, "application/json", bytes.NewBuffer(json_data))
				log.Println("Results sent")
			}()
			res.Body.Close()
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	lambda.Start(Listen)
}
