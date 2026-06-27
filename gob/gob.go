package gob

import (
	"encoding/gob"
	"os"
)

type Model struct {
	UserVectors  map[string][]float64
	MovieVectors map[string][]float64

	NumberOfLatentFactors   int
	LearningRate            float64
	RegularizationParameter float64
	NumberOfEpochs          int
}

func SaveModel(path string, model Model) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(model)
}

func LoadModel(path string) (Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return Model{}, err
	}
	defer file.Close()

	var model Model
	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&model)
	return model, err
}
