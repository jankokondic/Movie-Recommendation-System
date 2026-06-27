package model

import (
	"math"
	"root/constants"
)

type Metrics struct {
	MAE  float64
	RMSE float64
}

type Evaluation struct {
	MatrixFactorization Metrics
	GlobalAverage       Metrics
	UserAverage         Metrics
	MovieAverage        Metrics
}

func Test(newData []constants.Rating, model Model) Evaluation {
	var (
		mfAbs, mfSq         float64
		globalAbs, globalSq float64
		userAbs, userSq     float64
		movieAbs, movieSq   float64

		mfCount, globalCount, userCount, movieCount int
	)

	for _, value := range newData {
		// Matrix Factorization
		userVector, userExist := model.UserLatentFactor[value.UserID]
		movieVector, movieExist := model.MovieLatentFactor[value.MovieID]

		if userExist && movieExist {
			prediction := DotProduct(userVector, movieVector)
			err := value.Rating - prediction

			mfAbs += math.Abs(err)
			mfSq += err * err
			mfCount++
		}

		// Global Average
		{
			err := value.Rating - model.GlobalAverage
			globalAbs += math.Abs(err)
			globalSq += err * err
			globalCount++
		}

		// User Average
		if avg, exists := model.UserAverage[value.UserID]; exists {
			err := value.Rating - avg
			userAbs += math.Abs(err)
			userSq += err * err
			userCount++
		}

		// Movie Average
		if avg, exists := model.MovieAverage[value.MovieID]; exists {
			err := value.Rating - avg
			movieAbs += math.Abs(err)
			movieSq += err * err
			movieCount++
		}
	}

	return Evaluation{
		MatrixFactorization: Metrics{
			MAE:  mfAbs / float64(mfCount),
			RMSE: math.Sqrt(mfSq / float64(mfCount)),
		},
		GlobalAverage: Metrics{
			MAE:  globalAbs / float64(globalCount),
			RMSE: math.Sqrt(globalSq / float64(globalCount)),
		},
		UserAverage: Metrics{
			MAE:  userAbs / float64(userCount),
			RMSE: math.Sqrt(userSq / float64(userCount)),
		},
		MovieAverage: Metrics{
			MAE:  movieAbs / float64(movieCount),
			RMSE: math.Sqrt(movieSq / float64(movieCount)),
		},
	}
}
