package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Station struct {
	Name     string
	MinTemp  float64
	MaxTemp  float64
	MeanTemp float64
	count    int
}

func CalculateAverage(fileName string) {
	stationMap := map[string]*Station{}
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file: ", err)
		os.Exit(1)
	}
	defer f.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.IndexByte(line, ';')
		if idx <= 0 {
			continue
		}
		station, temperature := line[:idx], line[idx+1:]
		temp, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			log.Printf("error converting temperature to float: %v\n", err)
		}
		if val, ok := stationMap[station]; ok {
			if val.MaxTemp < temp {
				val.MaxTemp = temp
			}
			if val.MinTemp > temp {
				val.MinTemp = temp
			}
			mean := ((val.MeanTemp * float64(val.count)) + temp) / (float64(val.count) + 1)
			val.MeanTemp = math.Round(mean*10.0) / 10.0
			val.count++
		} else {
			stationMap[station] = &Station{Name: station, MinTemp: temp, MaxTemp: temp, MeanTemp: temp, count: 1}
		}
	}
	stationList := make([]Station, len(stationMap))
	var i int
	for _, station := range stationMap {
		stationList[i] = *station
		i++
	}

	slices.SortFunc(stationList, func(a, b Station) int {
		if a.Name > b.Name {
			return 1
		} else if a.Name < b.Name {
			return -1
		}
		return 0
	})
	fmt.Println(stationList)
}
