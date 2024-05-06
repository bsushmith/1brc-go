package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Station struct {
	Name    string
	MinTemp float64
	MaxTemp float64
	sum     float64
	count   int
}

func CalculateAverage(fileName string) {
	stationMap := map[string]*Station{}
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file: ", err)
		os.Exit(1)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	buf := make([]byte, 1<<20)
	var position int
	for {
		n, err := reader.Read(buf[position:])
		if err != nil && err != io.EOF {
			log.Fatal("error reading file: ", err)
		}
		position += n
		if err == io.EOF {
			break
		}

		for i := position - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				processLines(buf[:i+1], stationMap)
				copy(buf, buf[i+1:position])
				position -= i + 1
				break
			}
		}
	}
	if position > 0 {
		processLines(buf[:position], stationMap)
	}
	printStations(stationMap)
}

func printStations(stationMap map[string]*Station) {
	stationList := make([]string, len(stationMap))
	var i int
	for id, _ := range stationMap {
		stationList[i] = id
		i++
	}

	slices.Sort(stationList)
	fmt.Printf("{")
	for _, station := range stationList {
		st := stationMap[station]
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", st.Name, st.MinTemp, st.sum/float64(st.count), st.MaxTemp)
	}
	fmt.Printf("}\n")
}

func processLines(bytes []byte, stationMap map[string]*Station) {
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		idx := strings.IndexByte(line, ';')
		if idx < 0 {
			continue
		}
		station, temperature := line[:idx], line[idx+1:]
		temp, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			log.Printf("error converting temperature to float: %v\n", err)
		}
		if val, ok := stationMap[station]; ok {
			val.MaxTemp = max(val.MaxTemp, temp)
			val.MinTemp = min(val.MinTemp, temp)
			val.sum += temp
			val.count++
		} else {
			stationMap[station] = &Station{Name: station, MinTemp: temp, MaxTemp: temp, sum: temp, count: 1}
		}
	}
}
