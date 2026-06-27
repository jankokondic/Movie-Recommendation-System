package constants

const TrainFilePath string = "train.csv"
const TestFilePath string = "test.csv"
const Percentage int = 70
const ModelPath string = "model.j"

type Rating struct {
	UserID    int
	MovieID   int
	Rating    float64
	Timestamp string
}
