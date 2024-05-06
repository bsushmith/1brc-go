package main

import (
	"math"
	"math/rand/v2"
)

type WeatherStation struct {
	id              string
	meanTemperature float64
}

func (ws *WeatherStation) measurement() float64 {
	//sample = NormFloat64() * desiredStdDev + desiredMean
	temp := rand.NormFloat64()*10 + ws.meanTemperature
	return math.Round(temp*10.0) / 10.0
}
