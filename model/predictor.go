package model

import (
	"math"
	"root/constants"
)

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
