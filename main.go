package main

import (
	"fmt"
	"log"
	"root/constants"
	"root/model"
	"root/reader"
)

func main() {
	trainedModel, err := model.LoadModel(constants.ModelPath)
	if err != nil {
		log.Println(err)
		return
	}

	testData := reader.ReadTestData(constants.TestFilePath)
	evaluation := model.Test(testData, trainedModel)

	fmt.Println("MF:", evaluation.MatrixFactorization)
	fmt.Println("Global:", evaluation.GlobalAverage)
	fmt.Println("User:", evaluation.UserAverage)
	fmt.Println("Movie:", evaluation.MovieAverage)
}
