package main

import (
	"fmt"
	"log"
	"os"
	"root/constants"
	"root/model"
	"root/reader"
	"strconv"
)

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid %s=%q, using fallback %d", key, value, fallback)
		return fallback
	}

	return parsed
}

func getenvFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Invalid %s=%q, using fallback %f", key, value, fallback)
		return fallback
	}

	return parsed
}

func main() {
	conf := model.Configuration{
		NumberOfLatentFactors:   getenvInt("NUMBER_OF_LATENT_FACTORS", 20),
		LearningRate:            getenvFloat("LEARNING_RATE", 0.01),
		RegularizationParameter: getenvFloat("REGULARIZATION_PARAMETER", 0.02),
		NumberOfEpochs:          getenvInt("NUMBER_OF_EPOCHS", 20),
		InitializationMin:       getenvFloat("INITIALIZATION_MIN", -0.1),
		InitializationMax:       getenvFloat("INITIALIZATION_MAX", 0.1),
	}

	model.ModelRunner(conf)

	trainedModel, err := model.LoadModel(constants.ModelPath)
	if err != nil {
		log.Println(err)
		return
	}

	testData := reader.ReadTestData(constants.TestFilePath)

	loadMovie, err := reader.LoadMovies(constants.MoviePath)
	if err != nil {
		log.Println(err)
		return
	}

	listMovies := model.TopMoviesForUserByID(testData, trainedModel, 1, 10)
	for _, value := range listMovies {
		fmt.Println(loadMovie[value.MovieID])
	}

	evaluation := model.Test(testData, trainedModel)

	fmt.Println("MF:", evaluation.MatrixFactorization)
	fmt.Println("Global:", evaluation.GlobalAverage)
	fmt.Println("User:", evaluation.UserAverage)
	fmt.Println("Movie:", evaluation.MovieAverage)
	fmt.Println("Histogram", evaluation.ErrorHistogram)
}
