package statistic

import "root/constants"

type Statistic struct {
	GlobalAverage float64
	UserAverage   map[int]float64
	MovieAverage  map[int]float64
}

func New() Statistic {
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
