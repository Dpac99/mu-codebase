package tasks

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type Request struct {
	ID   string                 `json:"id"`
	Args map[string]interface{} `json:"args"`
}

type CPURequest struct {
	N       int         `json:"n"`
	Vectors [][]float64 `json:"vectors"`
}

type SleepRequest struct {
	Duration int `json:"duration"`
}

func ExecuteTask(req *Request) ([]byte, error) {
	log.Println(req)
	switch req.ID {
	case "sleep":
		duration := req.Args["duration"].(int)
		ret := placeHolderSleep(int(duration))
		b, err := json.Marshal(ret)
		if err != nil {
			return nil, err
		}

		return b, nil

	// case "cpu":
	// var r CPURequest
	// err := json.Unmarshal(req.Args, &r)
	// if err != nil {
	// 	return nil, err
	// }
	// ret := CPUintensive(&r)
	// b, err := json.Marshal(ret)
	// if err != nil {
	// 	return nil, err
	// }

	// return b, nil
	default:
		return nil, errors.New("unrecognized request")
	}
}

func placeHolderSleep(duration int) int {

	time.Sleep(time.Duration(duration * int(time.Second)))

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
