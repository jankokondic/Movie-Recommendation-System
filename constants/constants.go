package constants

const TrainFilePath string = "train.csv"
const TestFilePath string = "test.csv"
const Percentage int = 70

type Rating struct {
	UserID    int
	MovieID   int
	Rating    float64
	Timestamp string
}
