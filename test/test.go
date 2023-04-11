package main

import (
	"fmt"
)

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

func main() {
	a := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}

	b := [][]float64{
		{2, 3},
		{6, 7},
		{10, 11},
		{14, 15},
	}

	c := multiplyMatrix(a, b)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			fmt.Print(c[i][j])
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
