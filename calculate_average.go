package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Station struct {
	Name    string
	MinTemp float64
	MaxTemp float64
	sum     float64
	count   int
}

func CalculateAverage(fileName string) {
	numWorkers := runtime.NumCPU()

	outChan := make(chan map[string]*Station, numWorkers)
	var wg sync.WaitGroup

	finalResults := map[string]*Station{}
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file: ", err)
		os.Exit(1)
	}
	defer f.Close()

	var cwg sync.WaitGroup
	cwg.Add(1)
	go combineOutput(outChan, finalResults, &cwg)

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
				chunk := make([]byte, len(buf[:i+1]))
				copy(chunk, buf[:i+1])
				wg.Add(1)
				go processChunks(chunk, outChan, &wg)
				copy(buf, buf[i+1:position])
				position -= (i + 1)
				break
			}
		}
	}
	if position > 0 {
		wg.Add(1)
		processChunks(buf[:position], outChan, &wg)
	}
	wg.Wait()
	close(outChan)
	cwg.Wait()

	prettyPrint(finalResults)
}

func combineOutput(outChan chan map[string]*Station, finalResults map[string]*Station, cwg *sync.WaitGroup) {
	defer cwg.Done()
	for output := range outChan {
		for name, station := range output {
			if val, ok := finalResults[name]; ok {
				val.MaxTemp = max(val.MaxTemp, station.MaxTemp)
				val.MinTemp = min(val.MinTemp, station.MinTemp)
				val.sum += station.sum
				val.count += station.count
			} else {
				finalResults[name] = station
			}
		}
	}
}

func prettyPrint(finalResults map[string]*Station) {
	stationList := make([]string, len(finalResults))
	var i int
	for name := range finalResults {
		stationList[i] = name
		i++
	}
	slices.Sort(stationList)

	fmt.Printf("{")
	for j, station := range stationList {
		if j > 0 {
			fmt.Printf(", ")
		}
		st := finalResults[station]
		fmt.Printf("%s=%.1f/%.1f/%.1f", st.Name, st.MinTemp, st.sum/float64(st.count), st.MaxTemp)
	}
	fmt.Printf("}\n")
}

func processChunks(bytes []byte, outChan chan map[string]*Station, wg *sync.WaitGroup) {
	defer wg.Done()
	stationMap := make(map[string]*Station)
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
	outChan <- stationMap
}
