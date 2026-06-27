package main

import (
	"log"
	"root/constants"
	"root/model"
)

func main() {
	model, err := model.LoadModel(constants.ModelPath)
	if err != nil {
		log.Println(err)
		return
	}

	_ = model
}
