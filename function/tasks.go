package tasks

import (
	"time"
)

func PlaceHolderSleep(duration int) int {

	time.Sleep(time.Duration(duration * int(time.Second)))

	return 1
}

func CPUintensive(n int, vectors [][]float64) (result float64) {
	result = 0
	k := len(vectors)
	for i := 0; i < k; i += 2 {
		result += dotProduct(vectors[i], vectors[i+1])
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
