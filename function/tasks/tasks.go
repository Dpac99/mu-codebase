package tasks

import (
	"encoding/json"
	"errors"
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

type MatrixRequest struct {
	A [][]float64 `json:"a"`
	B [][]float64 `json:"b"`
}

func ExecuteTask(req types.TaskRequest) (interface{}, error) {
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

	case "matrix":
		var r MatrixRequest
		jsonData, _ := json.Marshal(req.Args)
		json.Unmarshal(jsonData, &r)
		ret := multiplyMatrix(r.A, r.B)
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

func multiplyMatrix(a [][]float64, b [][]float64) [][]float64 {
	n_rows := len(a)
	n_cols := len(b[0])
	n_elems := len(b)
	c := make([][]float64, n_rows)
	for i := range c {
		c[i] = make([]float64, n_cols)
	}

	for i := 0; i < n_rows; i++ {
		for j := 0; j < n_cols; j++ {
			c[i][j] = 0
			for k := 0; k < n_elems; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return c
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
