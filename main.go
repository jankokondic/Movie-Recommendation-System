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
	mea, rmse := model.Test(testData, trainedModel)

	fmt.Println(mea, rmse)
}
