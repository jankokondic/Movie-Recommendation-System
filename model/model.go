package model

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
	LearningRate            float64
	RegularizationParameter float64
	NumberOfEpochs          int
	InitializationMin       float64
	InitializationMax       float64
}

type Engine struct {
	User   map[int][]float64
	Movies map[int][]float64
	Data   []constants.Rating
	Configuration
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func New(data []constants.Rating, configuration Configuration) *Engine {
	e := &Engine{
		User:          make(map[int][]float64),
		Movies:        make(map[int][]float64),
		Data:          data,
		Configuration: configuration,
	}

	e.Init()

	return e
}

func (e *Engine) CreateLatentFactor() []float64 {
	numberOfLatentFactor := e.NumberOfLatentFactors

	latentFactor := make([]float64, numberOfLatentFactor)

	for index := range numberOfLatentFactor {
		latentFactor[index] = RandomFloat(e.InitializationMin, e.InitializationMax)
	}

	return latentFactor
}

func ensureLatentFactor(store map[int][]float64, id int, create func() []float64) {
	if _, exists := store[id]; !exists {
		store[id] = create()
	}
}

func (e *Engine) InitialLatentFactor(movieID, userID int) {
	ensureLatentFactor(e.User, userID, e.CreateLatentFactor)
	ensureLatentFactor(e.Movies, movieID, e.CreateLatentFactor)
}

func (e *Engine) Init() {
	for _, record := range e.Data {
		e.InitialLatentFactor(record.MovieID, record.UserID)
	}
}

func DotProduct(userLatentFactor, movieLatentFactor []float64) float64 {
	var product float64

	for index := range userLatentFactor {
		product += userLatentFactor[index] * movieLatentFactor[index]
	}

	return product
}

func Error(predictedRating float64, realRating float64) float64 {
	return realRating - predictedRating
}

func UpdateLatentFactors(
	userLatentFactor []float64,
	movieLatentFactor []float64,
	learningRate float64,
	ratingError float64,
	regularizationParameter float64,
) {
	for i := range userLatentFactor {
		oldUser := userLatentFactor[i]
		oldMovie := movieLatentFactor[i]

		userLatentFactor[i] = oldUser + learningRate*(ratingError*oldMovie-regularizationParameter*oldUser)
		movieLatentFactor[i] = oldMovie + learningRate*(ratingError*oldUser-regularizationParameter*oldMovie)
	}
}

func (e *Engine) Run() {
	for epoch := 0; epoch < e.NumberOfEpochs; epoch++ {
		for _, row := range e.Data {
			userLatentFactor := e.User[row.UserID]
			movieLatentFactor := e.Movies[row.MovieID]

			predictedRating := DotProduct(userLatentFactor, movieLatentFactor)
			ratingError := Error(predictedRating, row.Rating)

			UpdateLatentFactors(
				userLatentFactor,
				movieLatentFactor,
				e.LearningRate,
				ratingError,
				e.RegularizationParameter,
			)
		}
	}
}

func ModelRunner() {
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
	engine.Run()

	SaveModel(constants.ModelPath, Model{
		engine.User,
		engine.Movies,
		engine.NumberOfLatentFactors,
		engine.LearningRate,
		engine.RegularizationParameter,
		engine.NumberOfEpochs,
	})
	// for _, data := range engine.Data[:100] {
	// 	userList := engine.User[data.UserID]
	// 	movieList := engine.Movies[data.MovieID]

	// 	dot := DotProduct(userList, movieList)

	// 	fmt.Printf("real rating %f | my prediction %f | %f \n", data.Rating, dot, Error(dot, data.Rating))
	// }
}
