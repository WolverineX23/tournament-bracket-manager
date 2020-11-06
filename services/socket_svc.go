package services

import "github.com/sirupsen/logrus"

type SocketService struct {
	Log *logrus.Entry
}

func NewSocketService(logger *logrus.Logger) *SocketService {
	return &SocketService{
		Log: logger.WithField("services", "Socket"),
	}
}
