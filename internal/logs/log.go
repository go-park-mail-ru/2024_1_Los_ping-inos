package logs

import (
	"github.com/sirupsen/logrus"
	"os"
)

const RequestID = "requestID"

var Log = logrus.New()

func InitLog() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
}
