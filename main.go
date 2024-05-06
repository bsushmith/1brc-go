package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Usage: 1brc create_measurements <number of records to create>")
		os.Exit(1)
	}
	fileName := "./measurements.csv"

	if os.Args[1] == "create_measurements" && len(os.Args) == 3 {
		size, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid value for <number of records to create>")
			fmt.Println("Usage: ./bin/1brc create_measurements <number of records to create>")
			os.Exit(1)
		}
		CreateMeasurements(fileName, size)
	}

	if os.Args[1] == "calculate_average" {

		f, err := os.Create("cpu_profile")
		if err != nil {
			log.Fatalf("error creating cpu_profile file: %v\n", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		start := time.Now()
		CalculateAverage(fileName)
		elapsed := time.Since(start)
		fmt.Printf("Calculated average in %s \n", elapsed)
	}
}
