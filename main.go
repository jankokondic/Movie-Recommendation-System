package main

import (
	"fmt"
	"log"
	"root/constants"
	"root/model"
	"root/reader"
)

func main() {
	// conf := model.Configuration{
	// 	NumberOfLatentFactors:   20,
	// 	LearningRate:            0.01,
	// 	RegularizationParameter: 0.02,
	// 	NumberOfEpochs:          20,
	// 	InitializationMin:       -0.1,
	// 	InitializationMax:       0.1,
	// }

	// model.ModelRunner(conf)

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
