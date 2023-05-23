package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	A [][]int64 `json:"a"`
	B [][]int64 `json:"b"`
}

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) ([][]int64, error) {
	log.Println(req.Body)
	var r Request
	json.Unmarshal([]byte(req.Body), &r)
	log.Println(r)
	return multiplyMatrix(r.A, r.B), nil
}

func multiplyMatrix(a [][]int64, b [][]int64) [][]int64 {
	n_rows := len(a)
	n_cols := len(b[0])
	n_elems := len(b)
	c := make([][]int64, n_rows)
	for i := range c {
		c[i] = make([]int64, n_cols)
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

func main() {
	lambda.Start(HandleRequest)
}
