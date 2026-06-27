package main

import (
	"encoding/csv"
	"io"
	"os"
	"root/constants"
	"strconv"
)

type Configuration struct {
	NumberOfLatentFactors   int
	LearningRate            int
	RegularizationParameter int
	NumberOfEpochs          int
}

type Engine struct {
	User   map[string][]float64
	Movies map[string][]float64
	Data   []constants.Rating
	Configuration
}

func New(data []constants.Rating) *Engine {
	return &Engine{
		User:   make(map[string][]float64),
		Movies: make(map[string][]float64),
		Data:   data,
	}
}

func (e *Engine) Init() {

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

	engine := New(rating)

}
