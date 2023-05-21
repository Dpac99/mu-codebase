package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func avg(perc []float64) (res float64) {
	k := len(perc)
	for i := 0; i < k; i++ {
		res += perc[i]
	}
	res /= float64(k)
	return
}

func wait() {
	time.Sleep(2 * time.Second)
	cores := runtime.NumCPU()
	percents, _ := cpu.Percent(0, true)
	fmt.Printf("Number of cores: %d\n", cores)
	fmt.Printf("Number of cores (cpu): %d\n", len(percents))
	k := avg(percents)
	fmt.Printf("Percent avg: %f\n", k)

}

func main() {
	go wait()
	for true {
	}
}
