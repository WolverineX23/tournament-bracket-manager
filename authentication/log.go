package authentication

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func ConfigureLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = 4 // Info

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	LOG_FILE_LOCATION, exists := os.LookupEnv("LOG_FILE_LOCATION")
	if !exists {
		log.Fatal("missing LOG_FILE_LOCATION environment variable")
	}
	logfile, err := os.OpenFile(LOG_FILE_LOCATION, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to open file for log")
	} else {
		log.Out = logfile
		log.Formatter = &logrus.JSONFormatter{}
	}

	return log
}
