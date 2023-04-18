package main

import (
	"fmt"
	"log"
	"os"
)

var countFile = "/tmp/count.txt"

func main() {
	count := 0
	var F *os.File
	defer F.Close()
	_, err := os.Stat(countFile)
	if !os.IsNotExist(err) {
		F, err := os.Open(countFile)
		if err != nil {
			log.Fatalf(err.Error())
		}
		_, err = fmt.Fscanf(F, "%d\n", &count)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	F, err = os.Create(countFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("Count = %d\n", count)

	count++
	fmt.Printf("Count = %d\n", count)

	_, err = fmt.Fprintf(F, "%d", count)
	if err != nil {
		log.Fatalf(err.Error())
	}

}
