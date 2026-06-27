package main

import (
	"encoding/csv"
	"io"
	"math/rand/v2"
	"os"
	"root/constants"
	"strconv"
)

type Configuration struct {
	NumberOfLatentFactors   int
	LearningRate            float32
	RegularizationParameter float32
	NumberOfEpochs          int
	InitializationMin       float64
	InitializationMax       float64
}

type Engine struct {
	User   map[string][]float64
	Movies map[string][]float64
	Data   []constants.Rating
	Configuration
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func New(data []constants.Rating, configuration Configuration) *Engine {
	return &Engine{
		User:          make(map[string][]float64),
		Movies:        make(map[string][]float64),
		Data:          data,
		Configuration: configuration,
	}
}

func (e *Engine) InitialLatentFactors() []float64 {
	numberOfLatentFactor := e.NumberOfLatentFactors

	latentFactor := make([]float64, numberOfLatentFactor)

	for index := 0; index < numberOfLatentFactor; index++ {
		latentFactor[index] = RandomFloat(e.InitializationMin, e.InitializationMax)
	}

	return latentFactor
}

func (e *Engine) Init() {
	for _, record := range e.Data {

	}
}

func main() {
	inputFile, err := os.Open("rating.csv")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	var rating []constants.Rating

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		userID, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}

		movieID, err := strconv.Atoi(record[1])
		if err != nil {
			panic(err)
		}

		ratingValue, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			panic(err)
		}

		rating = append(rating, constants.Rating{
			UserID:    userID,
			MovieID:   movieID,
			Rating:    ratingValue,
			Timestamp: record[3],
		})
	}

	conf := Configuration{
		NumberOfLatentFactors:   20,
		LearningRate:            0.01,
		RegularizationParameter: 0.02,
		NumberOfEpochs:          20,
		InitializationMin:       -0.1,
		InitializationMax:       0.1,
	}

	engine := New(rating, conf)
	engine.Init()

}
