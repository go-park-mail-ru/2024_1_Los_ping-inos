package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

const RequestID = "requestID"
const Logg = "logger"

type Log struct {
	Logger    *logrus.Logger
	RequestID int64
}

func InitLog() Log {
	logger := Log{Logger: logrus.New()}
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Logger.Out = file
	} else {
		logger.Logger.Info("Failed to log to file, using default stderr")
	}
	return logger
}
