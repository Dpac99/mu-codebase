package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"serverless/types"
	"time"
)

type CPURequest struct {
	N       int         `json:"n"`
	Vectors [][]float64 `json:"vectors"`
}

type SleepRequest struct {
	Duration int `json:"duration"`
}

func ExecuteTask(req types.TaskRequest) (interface{}, error) {
	log.Println("Executing request of type " + req.Type + " with args " + string(fmt.Sprintf("%v", req.Args)))
	switch req.Type {
	case "sleep":
		duration := req.Args["duration"].(float64)
		log.Println(duration)
		ret := placeHolderSleep(float64(duration))

		return ret, nil

	case "cpu":
		var r CPURequest
		jsonData, _ := json.Marshal(req.Args)
		json.Unmarshal(jsonData, &r)
		// r.N = int(req.Args["n"].(float64))
		// r.Vectors = req.Args["vectors"].([][]float64)
		ret := CPUintensive(&r)

		return ret, nil

	case "test":
		name := req.Args["name"].(string)
		ret := testfunction(name)
		log.Println(ret)

		return ret, nil

	default:
		return nil, errors.New("unrecognized request")
	}
}

func placeHolderSleep(duration float64) int {

	time.Sleep(time.Duration(duration * float64(time.Second)))

	return 1
}

func CPUintensive(req *CPURequest) (result float64) {
	result = 0
	for i := 0; i < req.N; i += 2 {
		result += dotProduct(req.Vectors[i], req.Vectors[i+1])
	}
	return
}

func dotProduct(a []float64, b []float64) (result float64) {
	k := len(a)
	for i := 0; i < k; i++ {
		result += a[i] * b[i]
	}
	return
}

func testfunction(name string) string {
	return "hello " + name + "!"
}
