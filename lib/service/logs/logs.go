package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const LOGS_DIR = "./logs"

func New() *logrus.Logger {
	filename := fmt.Sprintf("%s/%s.log", LOGS_DIR, time.Now().Format("2006-01-02_15-04-05"))
	if err := os.MkdirAll(LOGS_DIR, os.ModePerm); err != nil {
		log.Panic(err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	logger := logrus.New()
	logger.Out = io.MultiWriter(os.Stdout, file)
	return logger
}
