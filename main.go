package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Rating struct {
	UserID    int
	MovieID   int
	Rating    float64
	Timestamp string
}

func main() {
	file, err := os.Open("rating.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var ratings []Rating

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		userID, _ := strconv.Atoi(record[0])
		movieID, _ := strconv.Atoi(record[1])
		rating, _ := strconv.ParseFloat(record[2], 64)

		ratings = append(ratings, Rating{
			UserID:    userID,
			MovieID:   movieID,
			Rating:    rating,
			Timestamp: record[3],
		})
	}

	fmt.Println("Broj ocjena:", len(ratings))
	fmt.Println(ratings[0])
}
