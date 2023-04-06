package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	N       int         `json:"n"`
	Vectors [][]float64 `json:"vectors"`
}

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (float64, error) {
	log.Println(req.Body)
	var r Request
	json.Unmarshal([]byte(req.Body), &r)
	var result float64 = 0
	for i := 0; i < r.N; i += 2 {
		result += dotProduct(r.Vectors[i], r.Vectors[i+1])
	}
	return result, nil
}

func dotProduct(a []float64, b []float64) (result float64) {
	k := len(a)
	for i := 0; i < k; i++ {
		result += a[i] * b[i]
	}
	return
}

func main() {
	lambda.Start(HandleRequest)
}
