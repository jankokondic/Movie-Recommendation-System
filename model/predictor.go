package model

import (
	"math"
	"root/constants"
	"sort"
)

type Statistic struct {
	GlobalAverage float64
	UserAverage   map[int]float64
	MovieAverage  map[int]float64
}

func NewStat() Statistic {
	return Statistic{
		UserAverage:  make(map[int]float64),
		MovieAverage: make(map[int]float64),
	}
}

func (s *Statistic) Load(data []constants.Rating) {
	var globalSum float64

	userSum := make(map[int]float64)
	userCount := make(map[int]int)

	movieSum := make(map[int]float64)
	movieCount := make(map[int]int)

	for _, value := range data {
		globalSum += value.Rating

		userSum[value.UserID] += value.Rating
		userCount[value.UserID]++

		movieSum[value.MovieID] += value.Rating
		movieCount[value.MovieID]++
	}

	s.GlobalAverage = globalSum / float64(len(data))

	for userID, sum := range userSum {
		s.UserAverage[userID] = sum / float64(userCount[userID])
	}

	for movieID, sum := range movieSum {
		s.MovieAverage[movieID] = sum / float64(movieCount[movieID])
	}
}

type Metrics struct {
	MAE  float64
	RMSE float64
}

type ErrorHistogram struct {
	ZeroTo025    int
	From025To050 int
	From050To075 int
	From075To100 int
	Above100     int
}

type Evaluation struct {
	MatrixFactorization Metrics
	GlobalAverage       Metrics
	UserAverage         Metrics
	MovieAverage        Metrics
	ErrorHistogram      ErrorHistogram
}

func AddToHistogram(hist *ErrorHistogram, absErr float64) {
	switch {
	case absErr < 0.25:
		hist.ZeroTo025++
	case absErr < 0.50:
		hist.From025To050++
	case absErr < 0.75:
		hist.From050To075++
	case absErr < 1.00:
		hist.From075To100++
	default:
		hist.Above100++
	}
}

func SafeMetrics(absSum, sqSum float64, count int) Metrics {
	if count == 0 {
		return Metrics{}
	}

	return Metrics{
		MAE:  absSum / float64(count),
		RMSE: math.Sqrt(sqSum / float64(count)),
	}
}

func Test(newData []constants.Rating, model Model) Evaluation {
	var (
		mfAbs, mfSq         float64
		globalAbs, globalSq float64
		userAbs, userSq     float64
		movieAbs, movieSq   float64

		mfCount, globalCount, userCount, movieCount int

		hist ErrorHistogram
	)

	for _, value := range newData {
		userVector, userExist := model.UserLatentFactor[value.UserID]
		movieVector, movieExist := model.MovieLatentFactor[value.MovieID]

		if userExist && movieExist {
			prediction := DotProduct(userVector, movieVector)
			err := value.Rating - prediction
			absErr := math.Abs(err)

			mfAbs += absErr
			mfSq += err * err
			mfCount++

			AddToHistogram(&hist, absErr)
		}

		errGlobal := value.Rating - model.GlobalAverage
		globalAbs += math.Abs(errGlobal)
		globalSq += errGlobal * errGlobal
		globalCount++

		if avg, exists := model.UserAverage[value.UserID]; exists {
			err := value.Rating - avg

			userAbs += math.Abs(err)
			userSq += err * err
			userCount++
		}

		if avg, exists := model.MovieAverage[value.MovieID]; exists {
			err := value.Rating - avg

			movieAbs += math.Abs(err)
			movieSq += err * err
			movieCount++
		}
	}

	return Evaluation{
		MatrixFactorization: SafeMetrics(mfAbs, mfSq, mfCount),
		GlobalAverage:       SafeMetrics(globalAbs, globalSq, globalCount),
		UserAverage:         SafeMetrics(userAbs, userSq, userCount),
		MovieAverage:        SafeMetrics(movieAbs, movieSq, movieCount),
		ErrorHistogram:      hist,
	}
}

type Predictor struct {
	MovieID int
	Rating  float64
}

func TopMoviesForUserByID(newData []constants.Rating, model Model, userID, topRange int) []Predictor {
	var listOfMovies []Predictor

	userVector, userExist := model.UserLatentFactor[userID]
	if !userExist {
		return nil
	}

	addedMovies := make(map[int]struct{})

	for _, value := range newData {
		if _, ok := addedMovies[value.MovieID]; ok {
			continue
		}

		movieVector, movieExist := model.MovieLatentFactor[value.MovieID]
		if !movieExist {
			continue
		}

		prediction := DotProduct(userVector, movieVector)

		listOfMovies = append(listOfMovies, Predictor{
			MovieID: value.MovieID,
			Rating:  prediction,
		})

		addedMovies[value.MovieID] = struct{}{}
	}

	sort.Slice(listOfMovies, func(i, j int) bool {
		return listOfMovies[i].Rating > listOfMovies[j].Rating
	})

	if topRange > len(listOfMovies) {
		topRange = len(listOfMovies)
	}

	return listOfMovies[:topRange]
}
