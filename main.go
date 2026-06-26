package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

const TrainFilePath string = "train.csv"
const TestFilePath string = "test.csv"
const Percentage int = 70

type Rating struct {
	UserID    int
	MovieID   int
	Rating    float64
	Timestamp string
}

func main() {
	inputFile, err := os.Open("rating.csv")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	os.Remove(TrainFilePath)
	os.Remove(TestFilePath)

	reader := csv.NewReader(inputFile)

	trainFile, err := os.OpenFile(TrainFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer trainFile.Close()

	testFile, err := os.OpenFile(TestFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer testFile.Close()

	trainWriter := csv.NewWriter(trainFile)
	defer trainWriter.Flush()

	testWriter := csv.NewWriter(testFile)
	defer testWriter.Flush()

	var currentUserID int
	var currentUserRatings []Rating
	firstRow := true

	for {
		record, err := reader.Read()

		if err == io.EOF {
			if len(currentUserRatings) > 0 {
				writeSplit(currentUserRatings, trainWriter, testWriter)
			}
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

		rating := Rating{
			UserID:    userID,
			MovieID:   movieID,
			Rating:    ratingValue,
			Timestamp: record[3],
		}

		if firstRow {
			currentUserID = userID
			currentUserRatings = append(currentUserRatings, rating)
			firstRow = false
			continue
		}

		if userID != currentUserID {
			writeSplit(currentUserRatings, trainWriter, testWriter)

			currentUserID = userID
			currentUserRatings = []Rating{rating}
		} else {
			currentUserRatings = append(currentUserRatings, rating)
		}
	}
}

func writeSplit(ratings []Rating, trainWriter *csv.Writer, testWriter *csv.Writer) {
	length := len(ratings)

	if length == 0 {
		return
	}

	trainCount := length * Percentage / 100

	if trainCount == 0 {
		trainCount = 1
	}

	if trainCount == length && length > 1 {
		trainCount = length - 1
	}

	writeRatings(trainWriter, ratings[:trainCount])
	writeRatings(testWriter, ratings[trainCount:])
}

func writeRatings(writer *csv.Writer, ratings []Rating) {
	for _, r := range ratings {
		err := writer.Write([]string{
			strconv.Itoa(r.UserID),
			strconv.Itoa(r.MovieID),
			strconv.FormatFloat(r.Rating, 'f', -1, 64),
			r.Timestamp,
		})
		if err != nil {
			panic(err)
		}
	}
}
