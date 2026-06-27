package model

import (
	"math"
	"root/constants"
)

func Test(newData []constants.Rating, model Model) (mae float64, rmse float64) {
	var (
		absoluteErrorSum float64
		squaredErrorSum  float64
		count            int
	)

	for _, value := range newData {
		userVector, userExist := model.UserLatentFactor[value.UserID]
		movieVector, movieExist := model.MovieLatentFactor[value.MovieID]

		if !userExist || !movieExist {
			continue
		}

		prediction := DotProduct(userVector, movieVector)
		ratingError := value.Rating - prediction
		absoluteErrorSum += math.Abs(ratingError)
		squaredErrorSum += ratingError * ratingError
		count++
	}

	return absoluteErrorSum / float64(count), math.Sqrt(squaredErrorSum / float64(count))
}
