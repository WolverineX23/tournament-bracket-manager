package authentication

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func ConfigureLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = 4 // Info

	/*LOG_FILE_LOCATION, exists := os.LookupEnv("LOG_FILE_LOCATION")
	if !exists {
		log.Fatal("missing LOG_FILE_LOCATION environment variable")
	}*/

	LOG_FILE_LOCATION := os.Getenv("NEW_LOG_FILE_LOCATION")
	fmt.Println(LOG_FILE_LOCATION)

	logfile, err := os.OpenFile(LOG_FILE_LOCATION, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to open file for log")
	} else {
		log.Out = logfile
		log.Formatter = &logrus.JSONFormatter{}
	}

	return log
}
